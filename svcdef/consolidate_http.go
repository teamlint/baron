package svcdef

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	gogen "github.com/golang/protobuf/protoc-gen-go/generator"

	"github.com/teamlint/baron/svcdef/svcparse"
)

type optional interface {
	Optional() bool
}

func isOptionalError(err error) bool {
	opt, ok := errors.Cause(err).(optional)
	return ok && opt.Optional()
}

func isEOF(err error) bool {
	if errors.Cause(err) == io.EOF || errors.Cause(err) == io.ErrUnexpectedEOF {
		return true
	}
	return false
}

// consolidateHTTP accepts a SvcDef and the io.Readers for the proto files
// comprising the definition. It modifies the SvcDef so that HTTPBindings and
// their associated HTTPParamters are added to each ServiceMethod. After this,
// each `HTTPBinding` will have a populated list of all the http parameters
// that that binding requires, where that parameter should be located, and the
// type of each parameter.
func consolidateHTTP(sd *Svcdef, protoFiles map[string]io.Reader) error {
	for _, pfile := range protoFiles {
		lex := svcparse.NewSvcLexer(pfile)
		protosvc, err := svcparse.ParseService(lex)
		if err != nil {
			if isOptionalError(err) {
				log.Warnf("Parser found rpc method which lacks HTTP " +
					"annotations; this is allowed, but will result in HTTP " +
					"transport not being generated.")
				return nil
			} else if isEOF(err) {
				continue
			}

			return errors.Wrap(err, "error while parsing http options for the service definition")
		}
		if sd == nil {
			return errors.New("Svcdef is nil")
		}
		if protosvc == nil {
			return errors.New("protosvc is nil")
		}
		// log.Printf("[svcdef/consolidate_http.go][consolidateHTTP] sd.Service=%+v, methods(%v)\n", sd.Service.Name, len(sd.Service.Methods))
		// service methods
		if sd.Service != nil {
			for _, m := range sd.Service.Methods {
				log.Debugf("\t %v.%v,%v Message = %+v\n", sd.Service.Name, m.Name, m.RequestType.Name, m.RequestType.Message)
				for _, f := range m.RequestType.Message.Fields {
					log.Debugf("\t\t %v.%v,%v Message.Field = %+v\n", sd.Service.Name, m.Name, m.RequestType.Name, f)
				}
			}
		}
		log.Debugf("consolidateHTTP http.protosvc=%+v", protosvc)
		for _, hm := range protosvc.Methods {
			log.Debugf("\t http.%v RequestType=%v, ResponseType=%v, HTTPBindings=%+v\n",
				hm.Name, hm.RequestType, hm.ResponseType, hm.HTTPBindings)
			for _, b := range hm.HTTPBindings {
				log.Debugf("\t\t http.Binding.Description = %+v\n", b.Description)
				log.Debugf("\t\t http.Binding.CustomHTTPPattern = %+v\n", b.CustomHTTPPattern)
				for _, f := range b.Fields {
					log.Debugf("\t\t\t http.Binding.Field = %+v\n", f)
				}
			}
		}
		err = assembleHTTPParams(sd.Service, protosvc)
		if err != nil {
			return errors.Wrap(err, "while assembling HTTP parameters")
		}
	}
	return nil
}

// assembleHTTPParams will use the output of the service parser to create
// HTTPParams for each service RequestType field indicating that parameters
// location, and the field to which it refers.
func assembleHTTPParams(svc *Service, httpsvc *svcparse.Service) error {
	getMethNamed := func(name string) *ServiceMethod {
		if svc == nil {
			return nil
		}
		for _, m := range svc.Methods {
			// Have to CamelCase the data from the parser since it may be lowercase
			// while the name from the Go file will be CamelCased
			if m.Name == gogen.CamelCase(name) {
				return m
			}
		}
		return nil
	}

	// This logic has been broken out of the for loop below to flatten
	// this function and avoid difficult to read nesting
	createParams := func(meth *ServiceMethod, parsedbind *svcparse.HTTPBinding) {
		msg := meth.RequestType.Message
		bind := HTTPBinding{}
		bind.Verb, bind.Path = getVerb(parsedbind)

		var params []*HTTPParameter
		for _, field := range msg.Fields {
			newParam := &HTTPParameter{}
			newParam.Field = field
			newParam.Location = paramLocation(field, parsedbind)
			params = append(params, newParam)
		}
		bind.Params = params
		meth.Bindings = append(meth.Bindings, &bind)
	}

	// Iterate through every HTTPBinding on every ServiceMethod, and create the
	// HTTPParameters for that HTTPBinding.
	for _, hm := range httpsvc.Methods {
		m := getMethNamed(hm.Name)
		if m == nil {
			return fmt.Errorf("cannot not find service method named %q", hm.Name)
		}
		for _, hbind := range hm.HTTPBindings {
			createParams(m, hbind)
		}
	}
	return nil
}

// getVerb returns the verb of a svcparse.HTTPBinding. The verb is found by
// first checking if there's a 'customHTTPPattern' for a binding and using
// that. If there's no custom verb defined, then we search through the defined
// fields for a 'standard' field such as 'get', 'post', etc. If the binding
// does not contain a field with a verb, returns two empty strings.
func getVerb(binding *svcparse.HTTPBinding) (verb string, path string) {
	if binding.CustomHTTPPattern != nil {
		for _, field := range binding.CustomHTTPPattern {
			if field.Kind == "kind" {
				verb = field.Value
			} else if field.Kind == "path" {
				path = field.Value
			}
		}
		return verb, path
	}
	for _, field := range binding.Fields {
		switch field.Kind {
		case "get", "put", "post", "delete", "patch":
			return field.Kind, field.Value
		}
	}
	return "", ""
}

// paramLocation returns the location that a field would be found according to
// the rules of a given HTTPBinding.
func paramLocation(field *Field, binding *svcparse.HTTPBinding) string {
	pathParams := getPathParams(binding)
	for _, param := range pathParams {
		// Have to CamelCase the data from the parser since it may be lowercase
		// while the name from the Go file will be CamelCased
		if gogen.CamelCase(strings.Split(param, ".")[0]) == field.Name {
			return "path"
		}
	}
	for _, optField := range binding.Fields {
		if optField.Kind == "body" {
			if optField.Value == "*" {
				return "body"
			} else if optField.Value == field.Name {
				return "body"
				// Have to CamelCase the fields from the protobuf file, as they may
				// be lowercase while the name from the Go file will be CamelCased.
			} else if gogen.CamelCase(strings.Split(optField.Value, ".")[0]) == field.Name {
				return "body"
			}
		}
	}

	return "query"
}

// Returns a slice of strings containing all parameters in the path
func getPathParams(binding *svcparse.HTTPBinding) []string {
	_, path := getVerb(binding)
	findParams := regexp.MustCompile("{(.*?)}")
	removeBraces := regexp.MustCompile("{|}")
	params := findParams.FindAllString(path, -1)
	rv := []string{}
	for _, p := range params {
		rv = append(rv, removeBraces.ReplaceAllString(p, ""))
	}
	return rv
}
