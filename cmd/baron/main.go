package main

import (
	"fmt"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/teamlint/baron/config"
	"github.com/teamlint/baron/internal/execprotoc"
	"github.com/teamlint/baron/internal/parsesvcname"
	"github.com/teamlint/baron/internal/start"

	ggkconf "github.com/teamlint/baron/gengokit"
	gengokit "github.com/teamlint/baron/gengokit/generator"
	"github.com/teamlint/baron/svcdef"
)

const (
	handlersDirName = "service"
	cmdDirName      = "cmd"
	serverDirName   = "server"
)

var (
	svcOutFlag    = flag.StringP("svcout", "o", "", "Go package path where the generated Go service will be written. Trailing slash will create a NAME-service directory")
	verboseFlag   = flag.BoolP("verbose", "v", false, "Verbose output")
	helpFlag      = flag.BoolP("help", "h", false, "Print usage")
	startFlag     = flag.BoolP("start", "s", false, "Output a 'start.proto' protobuf file in ./")
	versionFlag   = flag.BoolP("version", "V", false, "Print version")
	clientFlag    = flag.BoolP("client", "c", false, "Generate NAME-service client")
	transportFlag = flag.StringP("transport", "t", "all", "Service transport protocol: [grpc|nats]")
	svcdefFlag    = flag.BoolP("svcdef", "d", false, "Print service definition")
)

var binName = filepath.Base(os.Args[0])

var (
	// version is compiled into baron with the flag
	// go install -ldflags "-X main.version=$SHA"
	version string
	// BuildDate is compiled into baron with the flag
	// go install -ldflags "-X main.date=$VERSION_DATE"
	date string
	// buildinfo
	buildinfo string
)

func init() {
	// If Version or VersionDate are not set, baron was not built with make
	if version == "" || date == "" {
		rebuild := promptNoMake()
		if !rebuild {
			os.Exit(1)
		}
		err := makeAndRunbaron(os.Args)
		if err != nil {
			log.Fatal(errors.Wrap(err, "please install baron with make manually"))
		}
		os.Exit(0)
	}

	buildinfo = fmt.Sprintf("version: %s", version)
	buildinfo = fmt.Sprintf("%s version date: %s", buildinfo, date)

	flag.Usage = func() {
		if buildinfo != "" && (*verboseFlag || *helpFlag) {
			fmt.Fprintf(os.Stderr, "%s (%s)\n", binName, strings.TrimSpace(buildinfo))
		}
		fmt.Fprintf(os.Stderr, "\nUsage: %s [options] <protofile>...\n", binName)
		fmt.Fprintf(os.Stderr, "\nGenerates %s(go-kit) services using proto3 and gRPC definitions.\n", binName)
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Fprintf(os.Stderr, "%s (%s)\n", binName, strings.TrimSpace(buildinfo))
		os.Exit(0)
	}
	if *startFlag {
		pkg := ""
		outDir := ""
		if len(flag.Args()) > 0 {
			pkg = flag.Args()[0]
			if len(flag.Args()) > 1 {
				outDir = flag.Args()[1]
			}
		}
		os.Exit(start.Do(pkg, outDir))
	}

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "%s: missing .proto file(s)\n", binName)
		flag.Usage()
		os.Exit(1)
	}

	initLog()

	cfg, err := parseInput()
	if err != nil {
		log.Fatal(errors.Wrap(err, "cannot parse input"))
	}

	// If there was no service found in parseInput, the rest can be omitted.
	if cfg == nil {
		return
	}

	sd, err := parseServiceDefinition(cfg)
	if err != nil {
		log.Fatal(errors.Wrap(err, "cannot parse input definition proto files"))
	}

	// generate *.pb.baron.go
	err = generateBaronCode(cfg, sd)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "cannot generate {{.sd.PkgName}}.pb.baron.go"))
	}

	genFiles, err := generateGoKitCode(cfg, sd)
	if err != nil {
		log.Fatal(errors.Wrap(err, "cannot generate service"))
	}

	for path, file := range genFiles {
		dst := filepath.Join(cfg.ServicePath, path)
		err := gengokit.WriteGenFile(file, dst)
		if err != nil {
			log.Fatal(errors.Wrap(err, "cannot to write output"))
		}
		log.Infof(">> %s", dst)
	}

	cleanupOldFiles(cfg.ServicePath, strings.ToLower(sd.Service.Name))
}

// parseInput constructs a *config.Config with all values needed to parse
// service definition files.
func parseInput() (*config.Config, error) {
	cfg := config.Config{
		GenClient: *clientFlag,
		Transport: *transportFlag,
		Svcdef:    *svcdefFlag,
	}

	// GOPATH
	cfg.GoPath = filepath.SplitList(os.Getenv("GOPATH"))
	if len(cfg.GoPath) == 0 {
		cfg.GoPath = filepath.SplitList(build.Default.GOPATH)
	}
	log.WithField("GOPATH", cfg.GoPath).Debug()

	// DefPaths
	var err error
	rawDefinitionPaths := flag.Args()
	cfg.DefPaths, err = cleanProtofilePath(rawDefinitionPaths)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse input arguments")
	}
	log.WithField("DefPaths", cfg.DefPaths).Debug()

	protoDir := filepath.Dir(cfg.DefPaths[0])
	p, err := packages.Load(nil, protoDir) // 获取导入包
	if err != nil || len(p) == 0 {
		return nil, errors.Wrap(err, "proto files not found in importable go package")
	}

	cfg.PBPackage = p[0].PkgPath
	cfg.PBPath = protoDir
	log.WithField("PB Package", cfg.PBPackage).Debug()
	log.WithField("PB Path", cfg.PBPath).Debug()

	if err := execprotoc.GeneratePBDotGo(cfg.DefPaths, cfg.GoPath, cfg.PBPath); err != nil {
		return nil, errors.Wrap(err, "cannot create .pb.go files")
	}
	// 输出文件
	for _, p := range cfg.DefPaths {
		log.Infof("-> %s", parsesvcname.GetPBFileName(p, cfg.PBPath))
		log.Infof("-> %s", parsesvcname.GetGRPCPBFileName(p, cfg.PBPath))
	}

	// Service Path
	// 生成代码在临时目录, 获取服务名称
	svcName, err := parsesvcname.FromPaths(cfg.GoPath, cfg.DefPaths)
	if err != nil {
		log.Warnf("No valid service is defined; exiting now: %v", err)
		return nil, nil
	}

	svcName = strings.ToLower(svcName)

	svcDirName := svcName + "-service"
	log.WithField("svcDirName", svcDirName).Debug()

	// svcPath := filepath.Join(filepath.Dir(cfg.DefPaths[0]), svcDirName)
	// svcOut 服务输出目录不与 proto 关联,默认当前目录
	svcPath := filepath.Join(".", svcDirName)

	if *svcOutFlag != "" {
		svcOut := *svcOutFlag
		log.WithField("svcOutFlag", svcOut).Debug()

		// If the package flag ends in a seperator, file will be "".
		_, file := filepath.Split(svcOut)
		seperator := file == ""
		log.WithField("seperator", seperator).Debug()

		svcPath, err = parseSVCOut(svcOut, cfg.GoPath[0])
		if err != nil {
			return nil, errors.Wrapf(err, "cannot parse svcout: %s", svcOut)
		}

		// Join the svcDirName as a svcout ending with `/` should create it
		if seperator {
			svcPath = filepath.Join(svcPath, svcDirName)
		}
	}

	// log.WithField("svcPath", svcPath).Debug()

	// Create svcPath for the case that it does not exist
	err = os.MkdirAll(svcPath, 0777)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create svcPath directory: %s", svcPath)
	}

	p, err = packages.Load(nil, svcPath)
	if err != nil || len(p) == 0 {
		return nil, errors.Wrap(err, "generated service not found in importable go package")
	}

	log.WithField("Service Packages", p).Debug()

	cfg.ServicePackage = p[0].PkgPath
	cfg.ServicePath = svcPath

	log.WithField("Service Package", cfg.ServicePackage).Debug()
	log.WithField("Service PkgName", p[0].Name).Debug()
	log.WithField("Service Path", cfg.ServicePath).Debug()

	// PrevGen
	cfg.PrevGen, err = readPreviousGeneration(cfg.ServicePath, svcName)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read previously generated files")
	}

	return &cfg, nil
}

// parseSVCOut 解析服务输出目录
// 如果是相对路径转化为绝对路径
// 否则使用 GOPATH + svcOut 路径
func parseSVCOut(svcOut string, GOPATH string) (string, error) {
	if build.IsLocalImport(svcOut) {
		return filepath.Abs(svcOut)
	}
	return filepath.Join(GOPATH, "src", svcOut), nil
}

// parseServiceDefinition 返回Svcdef, 包含所有服务必要的信息
func parseServiceDefinition(cfg *config.Config) (*svcdef.Svcdef, error) {
	protoDefPaths := cfg.DefPaths
	// Create the ServicePath so the .pb.go files may be place within it
	if cfg.PrevGen == nil {
		err := os.MkdirAll(cfg.ServicePath, 0777)
		if err != nil {
			return nil, errors.Wrap(err, "cannot create service directory")
		}
	}

	// Get path names of .pb.go|grpc.pb.go files
	pbgoPaths := []string{}
	for _, p := range protoDefPaths {
		pbgoPaths = append(pbgoPaths, parsesvcname.GetPBFileName(p, cfg.PBPath))     // pb.go
		pbgoPaths = append(pbgoPaths, parsesvcname.GetGRPCPBFileName(p, cfg.PBPath)) // grpc.pb.go
	}
	pbgoFiles, err := openFiles(pbgoPaths)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open all .pb.go files")
	}

	pbFiles, err := openFiles(protoDefPaths)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open all .proto files")
	}

	// Create the svcdef
	sd, err := svcdef.New(pbgoFiles, pbFiles)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create service definition; did you pass ALL the protobuf files to baron?")
	}
	// 是否打印服务定义
	if cfg.Svcdef {
		fmt.Printf("%# v", pretty.Formatter(sd))
	}

	return sd, nil
}

// generateGoKitCode returns a map[string]io.Reader that represents a gokit
// service
func generateGoKitCode(cfg *config.Config, sd *svcdef.Svcdef) (map[string]io.Reader, error) {
	conf := ggkconf.Config{
		PBPackage:     cfg.PBPackage,
		GoPackage:     cfg.ServicePackage,
		PreviousFiles: cfg.PrevGen,
		GenClient:     cfg.GenClient,
		Transport:     cfg.Transport,
		Version:       version,
		VersionDate:   date,
	}

	// generate go-kit service
	genGokitFiles, err := gengokit.GenerateGokit(sd, conf)
	if err != nil {
		return nil, errors.Wrap(err, "cannot generate gokit service")
	}

	return genGokitFiles, nil
}

// generateBaronCode
func generateBaronCode(cfg *config.Config, sd *svcdef.Svcdef) error {
	conf := ggkconf.Config{
		PBPackage:     cfg.PBPackage,
		GoPackage:     cfg.ServicePackage,
		PreviousFiles: cfg.PrevGen,
		GenClient:     cfg.GenClient,
		Transport:     cfg.Transport,
		Version:       version,
		VersionDate:   date,
	}

	protoDefPaths := cfg.DefPaths

	// Get path names of .pb.baron.go files
	for _, p := range protoDefPaths {
		base := filepath.Base(p)
		barename := strings.TrimSuffix(base, filepath.Ext(p))
		baronPath := filepath.Join(cfg.PBPath, barename+".pb.baron.go")

		// generate go-kit service
		err := gengokit.GenerateBaronFile(sd, conf, baronPath)
		if err != nil {
			return errors.Wrap(err, "cannot generate baron service")
		}
		log.Infof("-> %s", baronPath)
	}

	return nil
}

func openFiles(paths []string) (map[string]io.Reader, error) {
	rv := map[string]io.Reader{}
	for _, p := range paths {
		reader, err := os.Open(p)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot open file %q", p)
		}
		rv[p] = reader
	}
	return rv, nil
}

// cleanProtofilePath returns the absolute filepath of a group of files
// of the files, or an error if the files are not in the same directory
func cleanProtofilePath(rawPaths []string) ([]string, error) {
	var fullPaths []string

	// Parsed passed file paths
	for _, def := range rawPaths {
		log.WithField("rawDefPath", def).Debug()
		full, err := filepath.Abs(def)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get working directory of baron")
		}
		log.WithField("fullDefPath", full)

		fullPaths = append(fullPaths, full)

		if filepath.Dir(fullPaths[0]) != filepath.Dir(full) {
			return nil, errors.Errorf("passed .proto files in different directories")
		}
	}

	return fullPaths, nil
}

// readPreviousGeneration returns a map[string]io.Reader representing the files in serviceDir
func readPreviousGeneration(serviceDir string, svcName string) (map[string]io.Reader, error) {
	if !fileExists(serviceDir) {
		return nil, nil
	}

	files := make(map[string]io.Reader)

	addFileToFiles := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			switch info.Name() {
			// Only files within the handlers dir are used to
			// support regeneration.
			// See `gengokit/generator/gen.go:generateResponseFile`
			case filepath.Base(serviceDir), svcName, svcName + "-client", handlersDirName, cmdDirName, serverDirName:
				return nil
			default:
				return filepath.SkipDir
			}
		}

		file, ioErr := os.Open(path)
		if ioErr != nil {
			return errors.Wrapf(ioErr, "cannot read file: %v", path)
		}

		// trim the prefix of the path to the proto files from the full path to the file
		relPath, err := filepath.Rel(serviceDir, path)
		if err != nil {
			return err
		}

		// ensure relPath is unix-style, so it matches what we look for later
		relPath = filepath.ToSlash(relPath)
		log.Infof("*> %s", filepath.Join(serviceDir, relPath))
		files[relPath] = file

		return nil
	}

	err := filepath.Walk(serviceDir, addFileToFiles)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot fully walk directory %v", serviceDir)
	}

	return files, nil
}

// fileExists checks if a file at the given path exists. Returns true if the
// file exists, and false if the file does not exist.
func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func cleanupOldFiles(servicePath, serviceName string) {
	serverCLI := filepath.Join(servicePath, "svc/server/cli")
	if _, err := os.Stat(serverCLI); err == nil {
		log.Warnf("Removing stale 'svc/server/cli' files")
		err := os.RemoveAll(serverCLI)
		if err != nil {
			log.Error(err)
		}
	}
	clientCLI := filepath.Join(servicePath, "svc/client/cli")
	if _, err := os.Stat(clientCLI); err == nil {
		log.Warnf("Removing stale 'svc/client/cli' files")
		err := os.RemoveAll(clientCLI)
		if err != nil {
			log.Error(err)
		}
	}

	oldServer := filepath.Join(servicePath, fmt.Sprintf("cmd/%s-server", serviceName))
	if _, err := os.Stat(oldServer); err == nil {
		log.Warnf(fmt.Sprintf("Removing stale 'cmd/%s-server' files, use cmd/%s going forward", serviceName, serviceName))
		err := os.RemoveAll(oldServer)
		if err != nil {
			log.Error(err)
		}
	}
}

// promptNoMake prints that baron was not built with make and prompts the user
// asking if they would like for this process to be automated
// returns true if yes, false if not.
func promptNoMake() bool {
	const msg = `
baron was not built using Makefile.
Please run 'make' inside go import path %s.

Do you want to automatically run 'make' and rerun command:

	$ `
	fmt.Printf(msg, baronImportPath)
	for _, a := range os.Args {
		fmt.Print(a)
		fmt.Print(" ")
	}
	const q = `

? [Y/n] `
	fmt.Print(q)

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	switch strings.ToLower(strings.TrimSpace(response)) {
	case "y", "yes":
		return true
	}
	return false
}

const baronImportPath = "github.com/teamlint/baron"

// makeAndRunbaron installs baron by running make in baronImportPath.
// It then passes through args to newly installed baron.
func makeAndRunbaron(args []string) error {
	p, err := build.Default.Import(baronImportPath, "", build.FindOnly)
	if err != nil {
		return errors.Wrap(err, "could not find baron directory")
	}
	make := exec.Command("make")
	make.Dir = p.Dir
	err = make.Run()
	if err != nil {
		return errors.Wrap(err, "could not run make in baron directory")
	}
	baron := exec.Command("baron", args[1:]...)
	baron.Stdin, baron.Stdout, baron.Stderr = os.Stdin, os.Stdout, os.Stderr
	return baron.Run()
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          false,
		DisableLevelTruncation: false,
		TimestampFormat:        "2006-01-02 15:04:05",
	})
	log.SetLevel(log.InfoLevel)
	if *verboseFlag {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}
}
