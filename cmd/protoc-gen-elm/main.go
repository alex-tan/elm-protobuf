package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jalandis/elm-protobuf/pkg/stringextras"

	"github.com/jalandis/elm-protobuf/pkg/elm"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "0.0.2"
const docUrl = "https://github.com/jalandis/elm-protobuf"

var excludedFiles = map[string]bool{
	"google/protobuf/timestamp.proto": true,
	"google/protobuf/wrappers.proto":  true,
}

type parameters struct {
	Version          bool
	Debug            bool
	RemoveDeprecated bool
}

func parseParameters(input *string) (parameters, error) {
	var result parameters
	var err error

	if input == nil {
		return result, nil
	}

	for _, i := range strings.Split(*input, ",") {
		switch i {
		case "remove-deprecated":
			result.RemoveDeprecated = true
		case "debug":
			result.Debug = true
		default:
			err = fmt.Errorf("unknown parameter: \"%s\"", i)
		}
	}

	return result, err
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Fprintf(os.Stdout, "%v %v\n", filepath.Base(os.Args[0]), version)
		os.Exit(0)
	}
	if len(os.Args) == 2 && os.Args[1] == "--help" {
		fmt.Fprintf(os.Stdout, "See "+docUrl+" for usage information.\n")
		os.Exit(0)
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Could not read request from STDIN: %v", err)
	}

	req := &pluginpb.CodeGeneratorRequest{}

	err = proto.Unmarshal(data, req)
	if err != nil {
		log.Fatalf("Could not unmarshal request: %v", err)
	}

	parameters, err := parseParameters(req.Parameter)
	if err != nil {
		log.Fatalf("Failed to parse parameters: %v", err)
	}

	if parameters.Debug {
		// Remove useless source code data.
		for _, inFile := range req.GetProtoFile() {
			inFile.SourceCodeInfo = nil
		}

		result, err := proto.Marshal(req)
		if err != nil {
			log.Fatalf("Failed to marshal request: %v", err)
		}

		log.Printf("Input data: %s", result)
	}

	plugins := (uint64)(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	resp := &pluginpb.CodeGeneratorResponse{
		SupportedFeatures: &plugins,
	}
	for _, inFile := range req.GetProtoFile() {
		log.Printf("Processing file %s", inFile.GetName())
		// Well Known Types.
		if excludedFiles[inFile.GetName()] {
			log.Printf("Skipping well known type")
			continue
		}

		name := fileName(inFile.GetName())
		content, err := templateFile(inFile, parameters)
		if err != nil {
			log.Fatalf("Could not template file: %v", err)
		}

		resp.File = append(resp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    &name,
			Content: &content,
		})
	}

	data, err = proto.Marshal(resp)
	if err != nil {
		log.Fatalf("Could not marshal response: %v [%v]", err, resp)
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		log.Fatalf("Could not write response to STDOUT: %v", err)
	}
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

func templateFile(inFile *descriptorpb.FileDescriptorProto, p parameters) (string, error) {
	t := template.New("t")

	t, err := elm.EnumCustomTypeTemplate(t)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse enum custom type template")
	}

	t, err = elm.OneOfCustomTypeTemplate(t)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse one-of custom type template")
	}

	t, err = elm.TypeAliasTemplate(t)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse type alias template")
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
		return "", errors.Wrap(err, "failed to parse nested PB message template")
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
{{- range .AdditionalImports }}
import {{ . }} exposing (..)
{{ end }}


uselessDeclarationToPreventErrorDueToEmptyOutputFile = 42
{{- range .TopEnums }}


{{ template "enum-custom-type" . }}
{{- end }}
{{- range .Messages }}


{{ template "nested-message" . }}
{{- end }}
`)
	if err != nil {
		return "", err
	}

	buff := &bytes.Buffer{}
	if err = t.Execute(buff, struct {
		SourceFile        string
		ModuleName        string
		ImportDict        bool
		AdditionalImports []string
		TopEnums          []elm.EnumCustomType
		Messages          []pbMessage
	}{
		SourceFile:        inFile.GetName(),
		ModuleName:        moduleName(inFile.GetName()),
		ImportDict:        hasMapEntries(inFile),
		AdditionalImports: getAdditionalImports(inFile.GetDependency()),
		TopEnums:          enumsToCustomTypes([]string{}, inFile.GetEnumType(), p),
		Messages:          messages([]string{}, inFile.GetMessageType(), p),
	}); err != nil {
		return "", err
	}

	return buff.String(), nil
}

type pbMessage struct {
	TypeAlias        elm.TypeAlias
	OneOfCustomTypes []elm.OneOfCustomType
	EnumCustomTypes  []elm.EnumCustomType
	NestedMessages   []pbMessage
}

func isDeprecated(options interface{}) bool {
	switch v := options.(type) {
	case *descriptorpb.MessageOptions:
		return v != nil && v.Deprecated != nil && *v.Deprecated
	case *descriptorpb.FieldOptions:
		return v != nil && v.Deprecated != nil && *v.Deprecated
	case *descriptorpb.EnumOptions:
		return v != nil && v.Deprecated != nil && *v.Deprecated
	case *descriptorpb.EnumValueOptions:
		return v != nil && v.Deprecated != nil && *v.Deprecated
	default:
		return false
	}
}

func enumsToCustomTypes(preface []string, enumPbs []*descriptorpb.EnumDescriptorProto, p parameters) []elm.EnumCustomType {
	var result []elm.EnumCustomType
	for _, enumPb := range enumPbs {
		if isDeprecated(enumPb.Options) && p.RemoveDeprecated {
			continue
		}

		var values []elm.EnumVariant
		for _, value := range enumPb.GetValue() {
			if isDeprecated(value.Options) && p.RemoveDeprecated {
				continue
			}

			values = append(values, elm.EnumVariant{
				Name:     elm.NestedVariantName(value.GetName(), preface),
				Number:   elm.ProtobufFieldNumber(value.GetNumber()),
				JSONName: elm.EnumVariantJSONName(value),
			})
		}

		enumType := elm.NestedType(enumPb.GetName(), preface)

		result = append(result, elm.EnumCustomType{
			Name:                   enumType,
			Decoder:                elm.DecoderName(enumType),
			Encoder:                elm.EncoderName(enumType),
			DefaultVariantVariable: elm.EnumDefaultVariantVariableName(enumType),
			DefaultVariantValue:    values[0].Name,
			Variants:               values,
		})
	}

	return result
}

func oneOfsToCustomTypes(preface []string, messagePb *descriptorpb.DescriptorProto, p parameters) []elm.OneOfCustomType {
	var result []elm.OneOfCustomType

	if isDeprecated(messagePb.Options) && p.RemoveDeprecated {
		return result
	}

	for oneofIndex, oneOfPb := range messagePb.GetOneofDecl() {
		var variants []elm.OneOfVariant
		for _, inField := range messagePb.GetField() {
			if isDeprecated(inField.Options) && p.RemoveDeprecated {
				continue
			}

			if inField.OneofIndex == nil || inField.GetOneofIndex() != int32(oneofIndex) {
				continue
			}

			variants = append(variants, elm.OneOfVariant{
				Name:     elm.NestedVariantName(inField.GetName(), preface),
				JSONName: elm.OneOfVariantJSONName(inField),
				Type:     elm.BasicFieldType(inField),
				Decoder:  elm.BasicFieldDecoder(inField),
				Encoder:  elm.BasicFieldEncoder(inField),
			})
		}

		name := elm.NestedType(oneOfPb.GetName(), preface)
		result = append(result, elm.OneOfCustomType{
			Name:     name,
			Decoder:  elm.DecoderName(name),
			Encoder:  elm.EncoderName(name),
			Variants: variants,
		})
	}

	return result
}

func proto3OptionalType(messagePb *descriptorpb.DescriptorProto, fieldPb *descriptorpb.FieldDescriptorProto) *elm.Type {
	oneofIndex := -1
	for i, v := range messagePb.GetOneofDecl() {
		if v.GetName() == fieldPb.GetTypeName() {
			oneofIndex = i
		}
	}
	if oneofIndex == -1 {
		fmt.Printf("no optional type found")
		return nil
	}
	fmt.Printf("found optional type for %s", fieldPb.GetTypeName())

	for _, inField := range messagePb.GetField() {
		if inField.GetProto3Optional() || inField.OneofIndex == nil || inField.GetOneofIndex() != int32(oneofIndex) {
			continue
		}
		v := elm.BasicFieldType(inField)
		return &v
	}
	return nil
}

func messages(preface []string, messagePbs []*descriptorpb.DescriptorProto, p parameters) []pbMessage {
	var result []pbMessage
	for _, messagePb := range messagePbs {
		if isDeprecated(messagePb.Options) && p.RemoveDeprecated {
			continue
		}

		var newFields []elm.TypeAliasField
		for _, fieldPb := range messagePb.GetField() {
			if isDeprecated(fieldPb.Options) && p.RemoveDeprecated {
				continue
			}

			if fieldPb.OneofIndex != nil {
				continue
			}

			nested := getNestedType(fieldPb, messagePb)
			proto3OptionalTypeResult := proto3OptionalType(messagePb, fieldPb)
			if proto3OptionalTypeResult != nil {
				newFields = append(newFields, elm.TypeAliasField{
					Name:    elm.FieldName(fieldPb.GetName()),
					Type:    elm.MaybeType(*proto3OptionalTypeResult),
					Number:  elm.ProtobufFieldNumber(fieldPb.GetNumber()),
					Encoder: elm.MapEncoder(fieldPb, nested),
					Decoder: elm.MapDecoder(fieldPb, nested),
				})
			} else if nested != nil {
				newFields = append(newFields, elm.TypeAliasField{
					Name:    elm.FieldName(fieldPb.GetName()),
					Type:    elm.MapType(nested),
					Number:  elm.ProtobufFieldNumber(fieldPb.GetNumber()),
					Encoder: elm.MapEncoder(fieldPb, nested),
					Decoder: elm.MapDecoder(fieldPb, nested),
				})
			} else if isOptional(fieldPb) {
				newFields = append(newFields, elm.TypeAliasField{
					Name:    elm.FieldName(fieldPb.GetName()),
					Type:    elm.MaybeType(elm.BasicFieldType(fieldPb)),
					Number:  elm.ProtobufFieldNumber(fieldPb.GetNumber()),
					Encoder: elm.MaybeEncoder(fieldPb),
					Decoder: elm.MaybeDecoder(fieldPb),
				})
			} else if isRepeated(fieldPb) {
				newFields = append(newFields, elm.TypeAliasField{
					Name:    elm.FieldName(fieldPb.GetName()),
					Type:    elm.ListType(elm.BasicFieldType(fieldPb)),
					Number:  elm.ProtobufFieldNumber(fieldPb.GetNumber()),
					Encoder: elm.ListEncoder(fieldPb),
					Decoder: elm.ListDecoder(fieldPb),
				})
			} else {
				newFields = append(newFields, elm.TypeAliasField{
					Name:    elm.FieldName(fieldPb.GetName()),
					Type:    elm.BasicFieldType(fieldPb),
					Number:  elm.ProtobufFieldNumber(fieldPb.GetNumber()),
					Encoder: elm.RequiredFieldEncoder(fieldPb),
					Decoder: elm.RequiredFieldDecoder(fieldPb),
				})
			}
		}

		for _, oneOfPb := range messagePb.GetOneofDecl() {
			newFields = append(newFields, elm.TypeAliasField{
				Name:    elm.FieldName(oneOfPb.GetName()),
				Type:    elm.OneOfType(oneOfPb.GetName()),
				Encoder: elm.OneOfEncoder(oneOfPb),
				Decoder: elm.OneOfDecoder(oneOfPb),
			})
		}

		newPreface := append([]string{messagePb.GetName()}, preface...)
		name := elm.NestedType(messagePb.GetName(), preface)
		result = append(result, pbMessage{
			TypeAlias: elm.TypeAlias{
				Name:    name,
				Decoder: elm.DecoderName(name),
				Encoder: elm.EncoderName(name),
				Fields:  newFields,
			},
			OneOfCustomTypes: oneOfsToCustomTypes([]string{}, messagePb, p),
			EnumCustomTypes:  enumsToCustomTypes(newPreface, messagePb.GetEnumType(), p),
			NestedMessages:   messages(newPreface, messagePb.GetNestedType(), p),
		})
	}

	return result
}

func isOptional(inField *descriptorpb.FieldDescriptorProto) bool {
	return inField.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL &&
		inField.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
}

func isRepeated(inField *descriptorpb.FieldDescriptorProto) bool {
	return inField.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED
}

func getLocalType(fullyQualifiedTypeName string) string {
	splitName := strings.Split(fullyQualifiedTypeName, ".")
	return splitName[len(splitName)-1]
}

func getNestedType(inField *descriptorpb.FieldDescriptorProto, inMessage *descriptorpb.DescriptorProto) *descriptorpb.DescriptorProto {
	localTypeName := getLocalType(inField.GetTypeName())
	for _, nested := range inMessage.GetNestedType() {
		if nested.GetName() == localTypeName && nested.GetOptions().GetMapEntry() {
			return nested
		}
	}

	return nil
}

func fileName(inFilePath string) string {
	inFileDir, inFileName := filepath.Split(inFilePath)

	trimmed := strings.TrimSuffix(inFileName, ".proto")
	shortFileName := stringextras.FirstUpper(trimmed)

	fullFileName := ""
	for _, segment := range strings.Split(inFileDir, "/") {
		if segment == "" {
			continue
		}

		fullFileName += stringextras.FirstUpper(segment) + "/"
	}

	return fullFileName + shortFileName + ".elm"
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
		if excludedFiles[d] {
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
