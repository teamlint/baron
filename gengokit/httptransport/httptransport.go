// Package httptransport provides functions and template helpers for templating
// the http-transport of a go-kit based service.
package httptransport

import (
	"bytes"
	"fmt"
	"go/format"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/teamlint/baron/pkg"
	"github.com/teamlint/baron/svcdef"
)

// Helper is the base struct for the data structure containing all the
// information necessary to correctly template the HTTP transport functionality
// of a service. Helper must be built from a Svcdef.
type Helper struct {
	Methods        []*Method
	ServerTemplate func(interface{}) (string, error)
	ClientTemplate func(interface{}) (string, error)
}

// NewHelper builds a helper struct from a service declaration. The other
// "New*" functions in this file are there to make this function smaller and
// more testable.
func NewHelper(svc *svcdef.Service) *Helper {
	// The HTTPAssistFuncs global is a group of function literals defined
	// within templates.go
	rv := Helper{
		// ServerTemplate: GenServerTemplate,
		// ClientTemplate: GenClientTemplate,
	}
	for _, meth := range svc.Methods {
		if len(meth.Bindings) > 0 {
			nMeth := NewMethod(meth)
			rv.Methods = append(rv.Methods, nMeth)
		}
	}
	return &rv
}

// NewMethod builds a Method struct from a svcdef.ServiceMethod.
func NewMethod(meth *svcdef.ServiceMethod) *Method {
	nMeth := Method{
		Name:         meth.Name,
		RequestType:  meth.RequestType.Name,
		ResponseType: meth.ResponseType.Name,
	}
	for i := range meth.Bindings {
		nBinding := NewBinding(i, meth)
		nBinding.Parent = &nMeth
		nMeth.Bindings = append(nMeth.Bindings, nBinding)
	}
	return &nMeth
}

// NewBinding 创建HTTPBinding, 通过索引,请求访问父级对象 SeeriveMethod, 创建方法绑定信息
func NewBinding(i int, meth *svcdef.ServiceMethod) *Binding {
	binding := meth.Bindings[i]
	nBinding := Binding{
		Label:        meth.Name + EnglishNumber(i),
		PathTemplate: binding.Path,
		BasePath:     basePath(binding.Path),
		Verb:         binding.Verb,
	}
	// Handle oneofs which need to be specially formed for query params
	for _, param := range binding.Params {
		// The 'Field' attr of each HTTPParameter always point to it's bound
		// Methods RequestType
		field := param.Field
		// only processing oneof fields
		if field.Type.Oneof == nil {
			continue
		}
		oneofField := OneofField{
			Name:     field.Name,
			Location: param.Location,
		}
		for _, oneofType := range field.Type.Oneof {
			option := Field{
				Name:           oneofType.Name,
				QueryParamName: oneofType.PBFieldName,
				CamelName:      pkg.GoCamelCase(field.Name),
				LowCamelName:   LowCamelName(oneofType.Name),
				Repeated:       oneofType.Type.ArrayType,
				GoType:         oneofType.Type.Name,
				LocalName:      fmt.Sprintf("%s%s", pkg.GoCamelCase(oneofType.Name), pkg.GoCamelCase(meth.Name)),
			}
			if oneofType.Type.Enum == nil && oneofType.Type.Map == nil {
				option.IsBaseType = true
			}
			// else {
			// 	option.GoType = "pb." + option.GoType
			// }

			// GoType 调整为基础类型,不包含 * 符号,* 号 通过 StarExpr 属性判断
			// 修改为指针类型或切片类型
			// if oneofType.Type.StarExpr {
			// 	option.GoType = "*" + option.GoType
			// if field.Type.ArrayType {
			// 	option.GoType = "[]" + option.GoType
			// }
			// }

			option.IsEnum = oneofType.Type.Enum != nil
			option.ConvertFunc, option.ConvertFuncNeedsErrorCheck = createDecodeConvertFunc(option)
			// option.TypeConversion = fmt.Sprintf("&pb.%s{%s: %s}", oneofType.Type.Message.Name, pkg.GoCamelCase(oneofType.Name), createDecodeTypeConversion(option))
			typeConvStr := "&%s{%s: %s}"
			// 如果是指针类型
			if oneofType.Type.StarExpr {
				typeConvStr = "%s{%s: %s}"
			}
			option.TypeConversion = fmt.Sprintf(typeConvStr, oneofType.Type.Message.Name, pkg.GoCamelCase(oneofType.Name), createDecodeTypeConversion(option))
			option.ZeroValue = getZeroValue(option)

			oneofField.Options = append(oneofField.Options, option)
		}
		nBinding.OneofFields = append(nBinding.OneofFields, &oneofField)
	}
	for _, param := range binding.Params {
		// The 'Field' attr of each HTTPParameter always point to it's bound
		// Methods RequestType
		field := param.Field
		// If the field is a oneof ignore; we handled above already
		if field.Type.Oneof != nil {
			continue
		}
		// 如果没有获取到有效字段类型,跳过
		if field.Type.Name == "" {
			continue
		}
		newField := Field{
			Name:           field.Name,
			QueryParamName: field.PBFieldName,
			CamelName:      pkg.GoCamelCase(field.Name),
			LowCamelName:   LowCamelName(field.Name),
			Location:       param.Location,
			Repeated:       field.Type.ArrayType,
			GoType:         field.Type.Name,
			StarExpr:       field.Type.StarExpr,
			LocalName:      fmt.Sprintf("%s%s", pkg.GoCamelCase(field.Name), pkg.GoCamelCase(meth.Name)),
		}
		// log.Debugf("NewBinding field = %+v\n", field.Type)

		if field.Type.Message == nil && field.Type.Enum == nil && field.Type.Map == nil && pkg.IsGoBaseType(newField.GoType) {
			newField.IsBaseType = true
		}
		// log.Debugf("NewBinding newField = %+v\n", newField)
		// else {
		// 	newField.GoType = "pb." + newField.GoType
		// }

		// 修改为指针类型或切片类型
		// GoType 调整为基础类型,不包含 * 符号,* 号 通过 StarExpr 属性判断
		// if field.Type.StarExpr {
		// 	newField.GoType = "*" + newField.GoType
		// 根据 Repeated 判断是否表添加 []
		// if field.Type.ArrayType {
		// 	newField.GoType = "[]" + newField.GoType
		// }
		// }

		// IsEnum needed for ConvertFunc and TypeConversion logic just below
		newField.IsEnum = field.Type.Enum != nil
		newField.ConvertFunc, newField.ConvertFuncNeedsErrorCheck = createDecodeConvertFunc(newField)
		newField.TypeConversion = createDecodeTypeConversion(newField)

		nBinding.Fields = append(nBinding.Fields, &newField)

		// Enums are allowed in query/path parameters, skip warning
		if newField.IsEnum {
			continue
		}

		// Emit warnings for certain cases
		if !newField.IsBaseType && newField.Location != "body" {
			log.Warnf(
				"%s.%s is a non-base type specified to be located outside of "+
					"the body. Non-base types outside the body may result in "+
					"generated code which fails to compile.",
				meth.Name,
				newField.Name)
		}
		if newField.Repeated && newField.Location == "path" {
			log.Warnf(
				"%s.%s is a repeated field specified to be in the path. "+
					"Repeated fields are not supported in the path and may"+
					"result in generated code which fails to compile.",
				meth.Name,
				newField.Name)
		}
	}
	return &nBinding
}

// func GenServerTemplate(exec interface{}) (string, error) {
// 	code, err := ApplyTemplate("ServerTemplate", templates.ServerTemplate, exec, TemplateFuncs)
// 	if err != nil {
// 		return "", err
// 	}
// 	code = FormatCode(code)
// 	return code, nil
// }

// func GenClientTemplate(exec interface{}) (string, error) {
// 	code, err := ApplyTemplate("ClientTemplate", templates.ClientTemplate, exec, TemplateFuncs)
// 	if err != nil {
// 		return "", err
// 	}
// 	code = FormatCode(code)
// 	return code, nil
// }

// GenServerDecode returns the generated code for the server-side decoding of
// an http request into its request struct.
// func (b *Binding) GenServerDecode() (string, error) {
// 	code, err := ApplyTemplate("ServerDecodeTemplate", templates.ServerDecodeTemplate, b, TemplateFuncs)
// 	if err != nil {
// 		return "", err
// 	}
// 	code = FormatCode(code)
// 	return code, nil
// }

// // GenClientEncode returns the generated code for the client-side encoding of
// // that clients request struct into the correctly formatted http request.
// func (b *Binding) GenClientEncode() (string, error) {
// 	code, err := ApplyTemplate("ClientEncodeTemplate", templates.ClientEncodeTemplate, b, TemplateFuncs)
// 	if err != nil {
// 		return "", err
// 	}
// 	code = FormatCode(code)
// 	return code, nil
// }

// PathSections returns a slice of strings for templating the creation of a
// fully assembled URL with the correct fields in the correct locations.
//
// For example, let's say there's a method "Sum" which accepts a "SumRequest",
// and SumRequest has two fields, 'a' and 'b'. Additionally, lets say that this
// binding for "Sum" has a path of "/sum/{a}". If we call the PathSection()
// method on this binding, it will return a slice that looks like the
// following slice literal:
//
//     []string{
//         "\"\"",
//         "\"sum\"",
//         "fmt.Sprint(req.A)",
//     }
func (b *Binding) PathSections() []string {
	isEnum := make(map[string]struct{})
	for _, v := range b.Fields {
		if v.IsEnum {
			isEnum[v.CamelName] = struct{}{}
		}
	}

	rv := []string{}
	parts := strings.Split(b.PathTemplate, "/")
	for _, part := range parts {
		if len(part) > 2 && part[0] == '{' && part[len(part)-1] == '}' {
			name := RemoveBraces(part)
			if _, ok := isEnum[pkg.GoCamelCase(name)]; ok {
				convert := fmt.Sprintf("fmt.Sprintf(\"%%d\", req.%v)", pkg.GoCamelCase(name))
				rv = append(rv, convert)
				continue
			}
			convert := fmt.Sprintf("fmt.Sprint(req.%v)", pkg.GoCamelCase(name))
			rv = append(rv, convert)
		} else {
			// Add quotes around things which'll be embeded as string literals,
			// so that the 'fmt.Sprint' lines will be unquoted and thus
			// evaluated as code.
			rv = append(rv, `"`+part+`"`)
		}
	}
	return rv
}

// GenQueryUnmarshaler returns the generated code for server-side unmarshaling
// of a query parameter into it's correct field on the request struct.
func (f *Field) GenQueryUnmarshaler() (string, error) {
	queryParamLogic := `
if {{.LocalName}}StrArr, ok := {{.Location}}Params["{{.QueryParamName}}"]; ok {
{{.LocalName}}Str := {{.LocalName}}StrArr[0]`

	pathParamLogic := `
{{.LocalName}}Str := {{.Location}}Params["{{.QueryParamName}}"]`

	genericLogic := `
{{.ConvertFunc}}{{if .ConvertFuncNeedsErrorCheck}}
if err != nil {
	return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting {{.LocalName}} from {{.Location}}, {{.Location}}Params: %v", {{.Location}}Params))
}{{end}}
req.{{.CamelName}} = {{.TypeConversion}}
`
	mergedLogic := queryParamLogic + genericLogic + "}"
	if f.Location == "path" {
		mergedLogic = pathParamLogic + genericLogic
	}

	code, err := ApplyTemplate("FieldEncodeLogic", mergedLogic, f, TemplateFuncs)
	if err != nil {
		return "", err
	}
	code = FormatCode(code)
	return code, nil
}

// GenQueryUnmarshaler returns the generated code for server-side unmarshaling
// of a query parameter into it's correct field on the request struct.
func (f *OneofField) GenQueryUnmarshaler() (string, error) {
	oneofEnclosure := `
{{- with $oneof := . -}}
	// {{.Name}} oneof
	{{.Name}}CountSet := 0
	{{range $option := $oneof.Options}}

		var {{$option.LocalName}}Str string
		{{$option.LocalName}}StrArr, {{$option.LocalName}}OK := {{$oneof.Location}}Params["{{$option.QueryParamName}}"]
		if {{$option.LocalName}}OK {
			{{$option.LocalName}}Str = {{$option.LocalName}}StrArr[0]
			{{$oneof.Name}}CountSet++
		}
	{{end}}

if {{.Name}}CountSet > 1 {
	return nil, errors.Errorf("only one of ({{range $option := $oneof.Options}}\"{{$option.QueryParamName}}\",{{end}}) allowed")
}

switch {
{{range $option := $oneof.Options}}
case {{$option.LocalName}}OK:
	{{$option.ConvertFunc}}{{if $option.ConvertFuncNeedsErrorCheck}}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting {{.LocalName}} from {{$option.Location}}, {{$oneof.Location}}Params: %v", {{$oneof.Location}}Params))
	}{{end}}
	req.{{$option.CamelName}} = {{$option.TypeConversion}}
{{end}}
}
{{- end -}}
`

	code, err := ApplyTemplate("FieldEncodeLogic", oneofEnclosure, f, TemplateFuncs)
	if err != nil {
		return "", err
	}
	code = FormatCode(code)
	return code, nil
}

// createDecodeConvertFunc creates a go string representing the function to
// convert the string form of the field to it's correct go type.
func createDecodeConvertFunc(f Field) (string, bool) {
	needsErrorCheck := true
	// log.Debugf("[createDecodeConvertFunc] field[%v].GoType=%v.IsBaseType=%v, Repeated=%v,filed.StarExpr=%v\n", f.Name, f.GoType, f.IsBaseType, f.Repeated, f.StarExpr)
	goType := f.GoType
	// 根据 Repeated 属性设置[]
	// goType := strings.TrimPrefix(f.GoType, "[]")
	// 指针类型处理,类型赋值时需要再处理
	// goType = strings.TrimPrefix(goType, "*")
	fType := ""
	switch goType {
	case "uint32":
		fType = "%s, err := strconv.ParseUint(%s, 10, 32)"
	case "uint64":
		fType = "%s, err := strconv.ParseUint(%s, 10, 64)"
	case "int32":
		fType = "%s, err := strconv.ParseInt(%s, 10, 32)"
	case "int64":
		fType = "%s, err := strconv.ParseInt(%s, 10, 64)"
	case "bool":
		fType = "%s, err := strconv.ParseBool(%s)"
	case "float32":
		fType = "%s, err := strconv.ParseFloat(%s, 32)"
	case "float64":
		fType = "%s, err := strconv.ParseFloat(%s, 64)"
	case "string":
		fType = "%s := %s"
		needsErrorCheck = false
	}

	if f.IsEnum && !f.Repeated {
		fType = "%s, err := strconv.ParseInt(%s, 10, 32)"
		return fmt.Sprintf(fType, f.LocalName, f.LocalName+"Str"), true
	}

	// 使用 json unmarshalling 任意自定义类型和重复类型
	if !f.IsBaseType || f.Repeated {
		// Args representing single custom message types are represented as
		// pointers. To do a bare assignment to a pointer, our rvalue must be a
		// pointer as well. So we special case args of a single custom message
		// type so that the variable LocalName is declared as a pointer.

		// var PageBaseStatsList *PageBase
		// PageBaseStatsList = &PageBase{}
		// err = json.Unmarshal([]byte(PageBaseStatsListStr), PageBaseStatsList)
		// if err != nil {
		// 	return nil, errors.Wrapf(err, "couldn't decode PageBaseStatsList from %v", PageBaseStatsListStr)
		// }
		// req.PageBase = PageBaseStatsList

		singleCustomTypeUnmarshalTmpl := `
var {{.LocalName}} {{if .StarExpr -}}*{{- end -}}{{.GoType}}
{{.LocalName}} = {{if .StarExpr -}}&{{- end -}}{{.GoType}}{}
err = json.Unmarshal([]byte({{.LocalName}}Str), {{.LocalName}})`

		errorCheckingTmpl := `
if err != nil {
	return nil, errors.Wrapf(err, "couldn't decode {{.LocalName}} from %v", {{.LocalName}}Str)
}`
		// All repeated args of any type are represented as slices, and bare
		// assignments to a slice accept a slice as the rvalue. As a result,
		// LocalName will be declared as a slice, and json.Unmarshal handles
		// everything else for us. Addititionally, if a type is a Base type and
		// is repeated, we first attempt to unmarshal the string we're
		// provided, and if that fails, we try to unmarshal the string
		// surrounded by square brackets. If THAT fails, then the string does
		// not represent a valid JSON string and an error is returned.

		// repeatedFieldType := strings.TrimPrefix(f.GoType, "[]")
		repeatedFieldType := f.GoType
		convertedVar := "converted"
		switch repeatedFieldType {
		case "uint32", "int32", "float32":
			convertedVar = repeatedFieldType + "(converted)"
		}
		repeatedUnmarshalTmpl := `
var {{.LocalName}} []{{if .StarExpr -}}*{{- end -}}{{.GoType}}
{{- if and (.IsBaseType) (not (Contains .GoType "[]byte"))}}
if len({{.LocalName}}StrArr) > 1 {
	{{- if (Contains .GoType "[]string")}}
	{{.LocalName}} = {{.LocalName}}StrArr
	{{- else}}
	{{.LocalName}} = make([]{{if .StarExpr -}}*{{- end -}}{{.GoType}}, len({{.LocalName}}StrArr))
	for i, v := range {{.LocalName}}StrArr {
	` + fmt.Sprintf(fType, "converted", "v") + errorCheckingTmpl + `
		{{.LocalName}}[i] = ` + convertedVar + `
	}
	{{- end}}
} else {
{{- end}}
	{{- if (Contains .GoType "[]string")}}
		{{.LocalName}} = strings.Split({{.LocalName}}Str, ",")
	{{- else if and (and .IsBaseType .Repeated) (not (Contains .GoType "[]byte"))}}
	err = json.Unmarshal([]byte({{.LocalName}}Str), &{{.LocalName}})
	if err != nil {
		{{.LocalName}}Str = "[" + {{.LocalName}}Str + "]"
	}
	err = json.Unmarshal([]byte({{.LocalName}}Str), &{{.LocalName}})
	{{- else}}
	err = json.Unmarshal([]byte({{.LocalName}}Str), &{{.LocalName}})
	{{- end}}
{{- if and (.IsBaseType) (not (Contains .GoType "[]byte"))}}
}
{{- end}}`

		var jsonConvTmpl string
		if !f.Repeated {
			jsonConvTmpl = singleCustomTypeUnmarshalTmpl + errorCheckingTmpl
		} else {
			jsonConvTmpl = repeatedUnmarshalTmpl
			if repeatedFieldType != "string" {
				jsonConvTmpl = repeatedUnmarshalTmpl + errorCheckingTmpl
			}
		}
		// log.Debugf("[UnmarshalNonBaseType] = %v", jsonConvTmpl)
		code, err := ApplyTemplate("UnmarshalNonBaseType", jsonConvTmpl, f, TemplateFuncs)
		if err != nil {
			panic(fmt.Sprintf("Couldn't apply template: %v", err))
		}
		return code, false
	}
	return fmt.Sprintf(fType, f.LocalName, f.LocalName+"Str"), needsErrorCheck
}

// createDecodeTypeConversion creates a go string that converts a 64 bit type
// to a 32 bit type as strconv.ParseInt, ParseUInt, and ParseFloat always
// return the 64 bit type. If the type is not a 64 bit integer type or is
// repeated, then returns the LocalName of that Field.
func createDecodeTypeConversion(f Field) string {
	if f.Repeated {
		// Equivalent of the 'default' case below, but taken early for repeated
		// types.
		return f.LocalName
	}
	fType := ""
	// log.Debugf("[gengokit/httptransport/httptransport.go][createDecodeTypeConversion] field=%+v\n", f)
	// goType := strings.TrimPrefix(f.GoType, "*")
	switch f.GoType {
	case "uint32", "int32", "float32":
		fType = f.GoType + "(%s)"
	default:
		fType = "%s"
	}
	if f.IsEnum {
		fType = f.GoType + "(%s)"
	}
	// 如果字段是标量类型并且是指针类型
	// log.Debugf("[createDecodeTypeConversion] field[%v].GoType=%v.IsBaseType=%v, filed.StarExpr=%v\n", f.Name, f.GoType, f.IsBaseType, f.StarExpr)
	if f.IsBaseType && f.StarExpr {
		return fmt.Sprintf("&"+fType, f.LocalName)
	}
	return fmt.Sprintf(fType, f.LocalName)
}

func getZeroValue(f Field) string {
	if !f.IsBaseType || f.Repeated || f.StarExpr {
		return "nil"
	}
	switch f.GoType {
	case "bool":
		return "false"
	case "string":
		return "\"\""
	default:
		return "0"
	}
}

// The 'basePath' of a path is the section from the start of the string till
// the first '{' character.
func basePath(path string) string {
	parts := strings.Split(path, "{")
	return parts[0]
}

// DigitEnglish is a map of runes of digits zero to nine to their lowercase
// english language spellings.
var DigitEnglish = map[rune]string{
	'0': "zero",
	'1': "one",
	'2': "two",
	'3': "three",
	'4': "four",
	'5': "five",
	'6': "six",
	'7': "seven",
	'8': "eight",
	'9': "nine",
}

// EnglishNumber takes an integer and returns the english words that represents
// that number, in base ten. Examples:
//     1  -> "One"
//     5  -> "Five"
//     10 -> "OneZero"
//     48 -> "FourEight"
func EnglishNumber(i int) string {
	n := strconv.Itoa(i)
	rv := ""
	for _, c := range n {
		if engl, ok := DigitEnglish[rune(c)]; ok {
			rv += strings.Title(engl)
		}
	}
	return rv
}

// LowCamelName returns a CamelCased string, but with the first letter
// lowercased. "example_name" becomes "exampleName".
func LowCamelName(s string) string {
	s = pkg.GoCamelCase(s)
	new := []rune(s)
	if len(new) < 1 {
		return s
	}
	rv := []rune{}
	rv = append(rv, unicode.ToLower(new[0]))
	rv = append(rv, new[1:]...)
	return string(rv)
}

// TemplateFuncs contains a series of utility functions to be passed into
// templates and used within those templates.
var TemplateFuncs = template.FuncMap{
	"ToLower":  strings.ToLower,
	"ToUpper":  strings.ToUpper,
	"Title":    strings.Title,
	"GoName":   pkg.GoCamelCase,
	"Contains": strings.Contains,
}

// ApplyTemplate applies a template with a given name, executor context, and
// function map. Returns the output of the template on success, returns an
// error if template failed to execute.
func ApplyTemplate(name string, tmpl string, data interface{}, fncs template.FuncMap) (string, error) {
	codeTemplate := template.Must(template.New(name).Funcs(fncs).Parse(tmpl))

	code := bytes.NewBuffer(nil)
	err := codeTemplate.Execute(code, data)
	if err != nil {
		return "", errors.Wrapf(err, "attempting to execute template %q", name)
	}
	return code.String(), nil
}

// FormatCode takes a string representing some go code and attempts to format
// that code. If formating fails, the original source code is returned.
func FormatCode(code string) string {
	formatted, err := format.Source([]byte(code))

	if err != nil {
		// Set formatted to code so at least we get something to examine
		formatted = []byte(code)
	}

	return string(formatted)
}
