package api

import (
	"bytes"
	"fmt"
	pgs "github.com/lyft/protoc-gen-star"
	"github.com/thematthopkins/elm-protobuf/pkg/elm"
	"github.com/thematthopkins/elm-protobuf/pkg/elmpb"
	"github.com/thematthopkins/elm-protobuf/pkg/stringextras"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"strings"
	"text/template"
)

func Generate(inFile *descriptorpb.FileDescriptorProto) (*pluginpb.CodeGeneratorResponse_File, error) {
	t := template.New("t")

	t, err := t.Parse(`module Api exposing (..)
-- DO NOT EDIT
-- AUTOGENERATED BY THE ELM PROTOCOL BUFFER COMPILER
-- https://github.com/tiziano88/elm-protobuf
-- source file: {{ .SourceFile }}

import Forward
import Helpers.Api
import Helpers.Api.Result exposing (ApiResult)
import HubTran.Effect as Effect
import HubTran.Flash as Flash
import Json.Decode
import Json.Encode
import Pb


defaultServerErrorHandler : error -> Effect.Effect msg
defaultServerErrorHandler error =
    let
        _ =
            Debug.log "Unexpected Result" error
    in
    Effect.flash Flash.defaultServerErrorAlert


request : String -> value -> (value -> Json.Encode.Value) -> Json.Decode.Decoder a -> Effect.Effect (ApiResult String a)
request path value encoder decoder =
    Helpers.Api.post path
        (encoder value)
        decoder
        |> Forward.sendWithError identity


{{ range .Endpoints }}
{{ .Name }} : {{.RequestType}} -> Effect.Effect (ApiResult String {{.ResponseType}})
{{ .Name }} p =
    request "{{.Path}}"
        p
        {{.Encoder}}
        {{.Decoder}}
{{ end }}
`)
	if err != nil {
		return nil, err
	}

	endpoints := []Endpoint{}

	packageName := stringextras.UpperCamelCase(elmpb.PackageName(inFile))

	services := inFile.GetService()
	for _, service := range services {
		for _, method := range service.Method {
			inputType := nameWithoutPackage(*method.InputType)
			outputType := nameWithoutPackage(*method.OutputType)
			serviceName := (pgs.Name)(*service.Name)
			methodName := (pgs.Name)(*method.Name)
			endpoints = append(endpoints, Endpoint{
				Name:         fmt.Sprintf("%s%s", stringextras.LowerCamelCase(*service.Name), stringextras.UpperCamelCase(*method.Name)),
				Path:         fmt.Sprintf("/%s/%s", serviceName.LowerSnakeCase(), methodName.LowerSnakeCase()),
				RequestType:  fmt.Sprintf("%s.%s", packageName, inputType),
				ResponseType: fmt.Sprintf("%s.%s", packageName, outputType),
				Encoder:      fmt.Sprintf("%s.%s", packageName, (string)(elm.EncoderName((elm.Type)(inputType)))),
				Decoder:      fmt.Sprintf("%s.%s", packageName, (string)(elm.DecoderName((elm.Type)(outputType)))),
			})
		}
	}

	buff := &bytes.Buffer{}
	if err = t.Execute(buff, struct {
		SourceFile string
		Endpoints  []Endpoint
	}{
		SourceFile: inFile.GetName(),
		Endpoints:  endpoints,
	}); err != nil {
		return nil, err
	}

	fileName := "Api.elm"
	result := buff.String()
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &fileName,
		Content: &result,
	}, nil
}

func nameWithoutPackage(v string) string {
	parts := strings.Split(v, ".")
	return parts[len(parts)-1]
}

type Endpoint struct {
	Name         string
	Path         string
	RequestType  string
	ResponseType string
	Encoder      string
	Decoder      string
}
