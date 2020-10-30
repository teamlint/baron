// Package generator generates a gokit service based on a service definition.
package generator

import (
	"bytes"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/teamlint/baron/gengokit"
	"github.com/teamlint/baron/gengokit/service"
	templFiles "github.com/teamlint/baron/gengokit/template"

	"github.com/teamlint/baron/svcdef"
)

var (
	ErrGenIgnored = errors.New("ignored files")
)

// GenerateGokit returns a gokit service generated from a service definition (svcdef),
// the package to the root of the generated service goPackage, the package
// to the .pb.go service struct files (goPBPackage) and any prevously generated files.
func GenerateGokit(sd *svcdef.Svcdef, conf gengokit.Config) (map[string]io.Reader, error) {
	data, err := gengokit.NewData(sd, conf)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create template data")
	}

	codeGenFiles := make(map[string]io.Reader)

	// Remove the suffix "-service" since it's added back in by templatePathToActual
	svcname := strings.ToLower(sd.Service.Name)
	for _, templPath := range templFiles.AssetNames() {

		// 排除不自动生成的模板
		if service.IsExcludedPath(templPath) {
			continue
		}
		// Re-derive the actual path for this file based on the service output
		// path provided by the baron main.go
		actualPath := templatePathToActual(templPath, svcname)
		file, err := generateResponseFile(templPath, data, conf.PreviousFiles[actualPath])
		if err != nil {
			if err == ErrGenIgnored {
				continue
			}
			return nil, errors.Wrap(err, "cannot render template")
		}
		codeGenFiles[actualPath] = file
	}

	return codeGenFiles, nil
}

// GenerateBaronFile returns a gokit service generated from a service definition (svcdef),
// the package to the root of the generated service goPackage, the package
// to the .pb.baron.go service struct files
func GenerateBaronFile(sd *svcdef.Svcdef, conf gengokit.Config, baronPath string) error {
	data, err := gengokit.NewData(sd, conf)
	if err != nil {
		return errors.Wrap(err, "cannot create baron template data")
	}

	templPath := service.BaronPath
	file, err := generateResponseFile(templPath, data, nil)
	if err != nil {
		return errors.Wrap(err, "cannot render baron template")
	}
	err = WriteGenFile(file, baronPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "cannot to write output"))
	}

	return nil
}

// generateResponseFile contains logic to choose how to render a template file
// based on path and if that file was generated previously. It accepts a
// template path to render, a templateExecutor to apply to the template, and a
// map of paths to files for the previous generation. It returns a
// io.Reader representing the generated file.
func generateResponseFile(templFP string, data *gengokit.Data, prevFile io.Reader) (io.Reader, error) {
	var genCode io.Reader
	var err error

	// Get the actual path to the file rather than the template file path
	actualFP := templatePathToActual(templFP, data.Service.Name)

	switch templFP {
	case service.ServicePath:
		h, err := service.NewService(data.Service, prevFile)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot parse previous handler: %q", actualFP)
		}

		if genCode, err = h.Render(templFP, data); err != nil {
			return nil, errors.Wrapf(err, "cannot render service template: %s", templFP)
		}
	case service.HookPath:
		hook := service.NewHook(prevFile)
		if genCode, err = hook.Render(templFP, data); err != nil {
			return nil, errors.Wrapf(err, "cannot render hook template: %s", templFP)
		}
	case service.MiddlewaresPath:
		m := service.NewMiddlewares(prevFile)
		// m.Load(prevFile)
		if genCode, err = m.Render(templFP, data); err != nil {
			return nil, errors.Wrapf(err, "cannot render middleware template: %s", templFP)
		}
	case service.CmdServerPath:
		r := service.NewCmdServer(prevFile)
		if genCode, err = r.Render(templFP, data); err != nil {
			return nil, errors.Wrapf(err, "cannot render cmd server template: %s", templFP)
		}
	case service.CmdClientPath:
		if data.Config.GenClient {
			r := service.NewCmdClient(prevFile)
			if genCode, err = r.Render(templFP, data); err != nil {
				return nil, errors.Wrapf(err, "cannot render cmd client template: %s", templFP)
			}
		} else {
			return nil, ErrGenIgnored
		}
	case service.ServerPath:
		r := service.NewServer(prevFile)
		if genCode, err = r.Render(templFP, data); err != nil {
			return nil, errors.Wrapf(err, "cannot render server template: %s", templFP)
		}
	default:
		if genCode, err = applyTemplateFromPath(templFP, data); err != nil {
			return nil, errors.Wrapf(err, "cannot render template: %s", templFP)
		}
	}

	codeBytes, err := ioutil.ReadAll(genCode)
	if err != nil {
		return nil, err
	}

	// ignore error as we want to write the code either way to inspect after
	// writing to disk
	formattedCode := formatCode(codeBytes)

	return bytes.NewReader(formattedCode), nil
}

// templatePathToActual accepts a templateFilePath and the svcName of the
// service and returns what the relative file path of what should be written to
// disk
func templatePathToActual(templFilePath, svcName string) string {
	// Switch "NAME" in path with svcName.
	// i.e. for svcName = addsvc; /NAME -> /addsvc-service/addsvc
	actual := strings.Replace(templFilePath, "NAME", svcName, -1)

	actual = strings.TrimSuffix(actual, ".tmpl")

	return actual
}

// applyTemplateFromPath calls applyTemplate with the template at templFilePath
func applyTemplateFromPath(templFP string, data *gengokit.Data) (io.Reader, error) {
	templBytes, err := templFiles.Asset(templFP)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find template file: %v", templFP)
	}

	return data.ApplyTemplate(string(templBytes), templFP)
}

// formatCode takes a string representing golang code and attempts to return a
// formated copy of that code.  If formatting fails, a warning is logged and
// the original code is returned.
func formatCode(code []byte) []byte {
	formatted, err := format.Source(code)

	if err != nil {
		log.WithError(err).Warn("Code formatting error, generated service will not build, outputting unformatted code")
		// return code so at least we get something to examine
		return code
	}

	return formatted
}

// WriteGenFile writes a file at path to the filesystem
func WriteGenFile(file io.Reader, path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "cannot create file %v", path)
	}

	_, err = io.Copy(outFile, file)
	if err != nil {
		return errors.Wrapf(err, "cannot write to %v", path)
	}
	return outFile.Close()
}
