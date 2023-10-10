package elmpb

import (
	"bytes"
	"google.golang.org/protobuf/types/pluginpb"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/thematthopkins/elm-protobuf/pkg/forwardextensions"
	"github.com/thematthopkins/elm-protobuf/pkg/generationparams"
	"github.com/thematthopkins/elm-protobuf/pkg/parsepb"
	"github.com/thematthopkins/elm-protobuf/pkg/stringextras"

	"github.com/thematthopkins/elm-protobuf/pkg/elm"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

var ExcludedFiles = map[string]bool{
	"forwardextensions.proto":          true,
	"google/protobuf/timestamp.proto":  true,
	"google/protobuf/wrappers.proto":   true,
	"google/protobuf/descriptor.proto": true,
}

func Generate(inFile *descriptorpb.FileDescriptorProto, p generationparams.Parameters) (*pluginpb.CodeGeneratorResponse_File, error) {
	t := template.New("t")

	t, err := elm.EnumCustomTypeTemplate(t)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse enum custom type template")
	}

	t, err = elm.OneOfCustomTypeTemplate(t)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse one-of custom type template")
	}

	t, err = elm.TypeAliasTemplate(t)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse type alias template")
	}

	t, err = t.Parse(`
{{- define "nested-message" -}}
{{ template "type-alias" .TypeAlias }}
{{- range .OneOfCustomTypes }}


{{ template "oneof-custom-type" . }}
{{- end }}
{{- range .EnumCustomTypes }}


{{ template "enum-custom-type" . }}
{{- end }}
{{- range .NestedMessages }}


{{ template "nested-message" . }}
{{- end }}
{{- end -}}
`)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse nested PB message template")
	}

	t, err = t.Parse(`module {{ .ModuleName }} exposing (..)

-- DO NOT EDIT
-- AUTOGENERATED BY THE ELM PROTOCOL BUFFER COMPILER
-- https://github.com/tiziano88/elm-protobuf
-- source file: {{ .SourceFile }}

import Protobuf exposing (..)

import Json.Decode as JD
import Json.Encode as JE
{{- if .ImportDict }}
import Dict
{{- end }}
{{- if .ImportIds }}
import Ids
{{- end }}
{{- range .AdditionalImports }}
import {{ . }} exposing (..)
{{ end }}


uselessDeclarationToPreventErrorDueToEmptyOutputFile = 42

requiredWithoutDefault : String -> JD.Decoder a -> JD.Decoder (a -> b) -> JD.Decoder b
requiredWithoutDefault name decoder d =
    field (JD.field name decoder) d

requiredFieldEncoderWithoutDefault : String -> (a -> JE.Value) -> a -> Maybe ( String, JE.Value )
requiredFieldEncoderWithoutDefault name encoder v =
    Just ( name, encoder v )

{{- range .TopEnums }}


{{ template "enum-custom-type" . }}
{{- end }}
{{- range .Messages }}


{{ template "nested-message" . }}
{{- end }}
`)
	if err != nil {
		return nil, err
	}

	buff := &bytes.Buffer{}
	if err = t.Execute(buff, struct {
		SourceFile        string
		ModuleName        string
		ImportDict        bool
		ImportIds         bool
		AdditionalImports []string
		TopEnums          []elm.EnumCustomType
		Messages          []parsepb.PbMessage
	}{
		SourceFile:        inFile.GetName(),
		ModuleName:        moduleName(PackageName(inFile)),
		ImportDict:        hasMapEntries(inFile),
		ImportIds:         p.GenerateForwardIds,
		AdditionalImports: getAdditionalImports(inFile.GetDependency()),
		TopEnums:          parsepb.EnumsToCustomTypes([]string{}, inFile.GetEnumType(), p),
		Messages:          parsepb.Messages([]string{}, inFile.GetMessageType(), p),
	}); err != nil {
		return nil, err
	}

	fileName := FileName(inFile)
	result := buff.String()
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &fileName,
		Content: &result,
	}, nil
}

func moduleName(inFilePath string) string {
	inFileDir, inFileName := filepath.Split(inFilePath)

	trimmed := strings.TrimSuffix(inFileName, ".proto")
	shortModuleName := stringextras.FirstUpper(trimmed)

	fullModuleName := ""
	for _, segment := range strings.Split(inFileDir, "/") {
		if segment == "" {
			continue
		}

		fullModuleName += stringextras.FirstUpper(segment) + "."
	}

	return fullModuleName + shortModuleName
}

func getAdditionalImports(dependencies []string) []string {
	var additions []string
	for _, d := range dependencies {
		if ExcludedFiles[d] {
			continue
		}

		fullModuleName := ""
		for _, segment := range strings.Split(strings.TrimSuffix(d, ".proto"), "/") {
			if segment == "" {
				continue
			}
			fullModuleName += stringextras.FirstUpper(segment) + "."
		}

		additions = append(additions, strings.TrimSuffix(fullModuleName, "."))
	}
	return additions
}

func PackageName(inFile *descriptorpb.FileDescriptorProto) string {
	if proto.HasExtension(inFile.Options, forwardextensions.E_ElmPackage) {
		return proto.GetExtension(inFile.Options, forwardextensions.E_ElmPackage).(string)
	}
	inFileDir, inFileName := filepath.Split(inFile.GetName())

	trimmed := strings.TrimSuffix(inFileName, ".proto")
	shortFileName := stringextras.FirstUpper(trimmed)

	fullFileName := ""
	for _, segment := range strings.Split(inFileDir, "/") {
		if segment == "" {
			continue
		}

		fullFileName += stringextras.FirstUpper(segment) + "/"
	}

	return fullFileName + shortFileName
}

func FileName(inFile *descriptorpb.FileDescriptorProto) string {
	return PackageName(inFile) + ".elm"
}

func hasMapEntries(inFile *descriptorpb.FileDescriptorProto) bool {
	for _, m := range inFile.GetMessageType() {
		if hasMapEntriesInMessage(m) {
			return true
		}
	}

	return false
}

func hasMapEntriesInMessage(inMessage *descriptorpb.DescriptorProto) bool {
	if inMessage.GetOptions().GetMapEntry() {
		return true
	}

	for _, m := range inMessage.GetNestedType() {
		if hasMapEntriesInMessage(m) {
			return true
		}
	}

	return false
}
