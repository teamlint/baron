package gengokit

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/teamlint/baron/gengokit/httptransport"
	templFiles "github.com/teamlint/baron/gengokit/template"
	"github.com/teamlint/baron/pkg"
	"github.com/teamlint/baron/svcdef"
)

// RenderStatus 表示 Renderable 状态的接口
type RenderStatus interface {
	IsFirst() bool    // 是否首次呈现
	IsModified() bool // 是否改写
}

type Renderable interface {
	Render(string, *Data) (io.Reader, error)
}

type Config struct {
	GoPackage   string
	ServicePath string // 服务输出目录
	PBPackage   string
	PBPath      string // .pb.go 输出目录
	Version     string
	VersionDate string

	PreviousFiles map[string]io.Reader
	// transport all|grpc|nats
	Transport string
	// generate client CLI
	GenClient bool
}

// FuncMap contains a series of utility functions to be passed into
// templates and used within those templates.
var FuncMap = template.FuncMap{
	"ToLower":  strings.ToLower,
	"ToUpper":  strings.ToUpper,
	"Title":    strings.Title,
	"GoName":   pkg.GoCamelCase,
	"Contains": strings.Contains,
}

// Data is passed to templates as the executing struct; its fields
// and methods are used to modify the template
type Data struct {
	// import path for the directory containing the definition .proto files
	ImportPath string
	// import path for .pb.go files containing service structs
	PBImportPath string
	// PackageName is the name of the package containing the service definition
	PackageName string
	// GRPC/Protobuff service, with all parameters and return values accessible
	Service *svcdef.Service
	// A helper struct for generating http transport functionality.
	HTTPHelper *httptransport.Helper
	FuncMap    template.FuncMap
	// generate config
	Config Config

	Version     string
	VersionDate string
}

// NewData GoKit 数据源
func NewData(sd *svcdef.Svcdef, conf Config) (*Data, error) {
	return &Data{
		ImportPath:   conf.GoPackage,
		PBImportPath: conf.PBPackage,
		PackageName:  sd.PkgName,
		Service:      sd.Service,
		HTTPHelper:   httptransport.NewHelper(sd.Service),
		FuncMap:      FuncMap,
		Version:      conf.Version,
		VersionDate:  conf.VersionDate,
		Config:       conf,
	}, nil
}

// ApplyTemplateFromPath 使用模板路径执行模板呈现
func (e *Data) ApplyTemplateFromPath(templFP string) (io.Reader, error) {
	tmplContent, err := templFiles.AssetString(templFP)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find template file: %v", templFP)
	}

	return e.ApplyTemplate(tmplContent, templFP)
}

// ApplyTemplate applies the passed template with the Data
func (e *Data) ApplyTemplate(templ string, templName string) (io.Reader, error) {
	return ApplyTemplate(templ, templName, e, e.FuncMap)
}

// ApplyTemplate is a helper methods that packages can call to render a
// template with any data and func map
func ApplyTemplate(templ string, templName string, data interface{}, funcMap template.FuncMap) (io.Reader, error) {
	codeTemplate, err := template.New(templName).Funcs(funcMap).Parse(templ)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create template")
	}

	outputBuffer := bytes.NewBuffer(nil)
	err = codeTemplate.Execute(outputBuffer, data)
	if err != nil {
		return nil, errors.Wrap(err, "template error")
	}

	return outputBuffer, nil
}
