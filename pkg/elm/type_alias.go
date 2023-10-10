package elm

import (
	"fmt"
	"text/template"

	"github.com/thematthopkins/elm-protobuf/pkg/stringextras"

	"google.golang.org/protobuf/types/descriptorpb"
)

// WellKnownType - information to handle Google well known types
type WellKnownType struct {
	Type    Type
	Encoder VariableName
	Decoder VariableName
}

var (
	// WellKnownTypeMap - map of Google well known type PB identifier to encoder/decoder info
	WellKnownTypeMap = map[string]WellKnownType{
		".google.protobuf.Timestamp": {
			Type:    "Timestamp",
			Decoder: "timestampDecoder",
			Encoder: "timestampEncoder",
		},
		".google.protobuf.Int32Value": {
			Type:    intType,
			Decoder: "intValueDecoder",
			Encoder: "intValueEncoder",
		},
		".google.protobuf.Int64Value": {
			Type:    intType,
			Decoder: "intValueDecoder",
			Encoder: "numericStringEncoder",
		},
		".google.protobuf.UInt32Value": {
			Type:    intType,
			Decoder: "intValueDecoder",
			Encoder: "intValueEncoder",
		},
		".google.protobuf.UInt64Value": {
			Type:    intType,
			Decoder: "intValueDecoder",
			Encoder: "numericStringEncoder",
		},
		".google.protobuf.DoubleValue": {
			Type:    floatType,
			Decoder: "floatValueDecoder",
			Encoder: "floatValueEncoder",
		},
		".google.protobuf.FloatValue": {
			Type:    floatType,
			Decoder: "floatValueDecoder",
			Encoder: "floatValueEncoder",
		},
		".google.protobuf.StringValue": {
			Type:    stringType,
			Decoder: "stringValueDecoder",
			Encoder: "stringValueEncoder",
		},
		".google.protobuf.BytesValue": {
			Type:    bytesType,
			Decoder: "bytesValueDecoder",
			Encoder: "bytesValueEncoder",
		},
		".google.protobuf.BoolValue": {
			Type:    boolType,
			Decoder: "boolValueDecoder",
			Encoder: "boolValueEncoder",
		},
	}

	reservedKeywords = map[string]bool{
		"module":   true,
		"exposing": true,
		"import":   true,
		"type":     true,
		"let":      true,
		"in":       true,
		"if":       true,
		"then":     true,
		"else":     true,
		"where":    true,
		"case":     true,
		"of":       true,
		"port":     true,
		"as":       true,
	}
)

// TypeAlias - defines an Elm type alias (somtimes called a record)
// https://guide.elm-lang.org/types/type_aliases.html
type TypeAlias struct {
	Name        Type
	IsSingleton bool
	LowerName   string
	Decoder     VariableName
	Encoder     VariableName
	Fields      []TypeAliasField
}

// FieldDecoder used in type alias decdoer (ex. )
type FieldDecoder string

// FieldEncoder used in type alias decdoer (ex. )
type FieldEncoder string

// Custom attribute to be able to force an id type
type IdTypeOverride string

// TypeAliasField - type alias field definition
type TypeAliasField struct {
	Name           VariableName
	Type           Type
	Number         ProtobufFieldNumber
	Decoder        FieldDecoder
	Encoder        FieldEncoder
	IdTypeOverride *IdTypeOverride
}

func appendUnderscoreToReservedKeywords(in string) string {
	if reservedKeywords[in] {
		return fmt.Sprintf("%s_", in)
	}

	return in
}

// FieldName - simple camelcase variable name with first letter lower
func FieldName(in string) VariableName {
	return VariableName(appendUnderscoreToReservedKeywords(stringextras.LowerCamelCase(in)))
}

// FieldJSONName - JSON identifier for field decoder/encoding
func FieldJSONName(pb *descriptorpb.FieldDescriptorProto) VariantJSONName {
	return VariantJSONName(pb.GetJsonName())
}

func RequiredFieldEncoder(parentName *string, pb *descriptorpb.FieldDescriptorProto) FieldEncoder {
	_, isWellKnownType := WellKnownTypeMap[pb.GetTypeName()]

	if pb.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE && !isWellKnownType {
		return FieldEncoder(fmt.Sprintf(
			"requiredFieldEncoderWithoutDefault \"%s\" %s v.%s",
			FieldJSONName(pb),
			BasicFieldEncoder(parentName, pb),
			FieldName(pb.GetName()),
		))
	}
	return FieldEncoder(fmt.Sprintf(
		"requiredFieldEncoder \"%s\" %s %s v.%s",
		FieldJSONName(pb),
		BasicFieldEncoder(parentName, pb),
		BasicFieldDefaultValue(parentName, pb),
		FieldName(pb.GetName()),
	))
}

func RequiredFieldDecoder(parentName *string, pb *descriptorpb.FieldDescriptorProto) FieldDecoder {
	_, isWellKnownType := WellKnownTypeMap[pb.GetTypeName()]

	if pb.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE && !isWellKnownType {
		return FieldDecoder(fmt.Sprintf(
			"requiredWithoutDefault \"%s\" %s",
			FieldJSONName(pb),
			BasicFieldDecoder(parentName, pb),
		))
	}
	return FieldDecoder(fmt.Sprintf(
		"required \"%s\" %s %s",
		FieldJSONName(pb),
		BasicFieldDecoder(parentName, pb),
		BasicFieldDefaultValue(parentName, pb),
	))
}

func OneOfEncoder(preface []string, pb *descriptorpb.OneofDescriptorProto) FieldEncoder {
	return FieldEncoder(fmt.Sprintf("%s v.%s",
		EncoderName(NestedType(stringextras.CamelCase(pb.GetName()), preface)),
		FieldName(pb.GetName()),
	))
}

func OneOfDecoder(preface []string, pb *descriptorpb.OneofDescriptorProto) FieldDecoder {
	return FieldDecoder(fmt.Sprintf(
		"field %s",
		DecoderName(NestedType(stringextras.CamelCase(pb.GetName()), preface)),
	))
}

func MapType(messagePb *descriptorpb.DescriptorProto) Type {
	keyField := messagePb.GetField()[0]
	valueField := messagePb.GetField()[1]

	return Type(fmt.Sprintf(
		"Dict.Dict %s %s",
		BasicFieldType(nil, keyField),
		BasicFieldType(nil, valueField),
	))
}

func MapEncoder(fieldPb *descriptorpb.FieldDescriptorProto,
	messagePb *descriptorpb.DescriptorProto,
) FieldEncoder {
	valueField := messagePb.GetField()[1]

	return FieldEncoder(fmt.Sprintf(
		"mapEntriesFieldEncoder \"%s\" %s v.%s",
		FieldJSONName(fieldPb),
		BasicFieldEncoder(nil, valueField),
		FieldName(fieldPb.GetName()),
	))
}

func MapDecoder(
	fieldPb *descriptorpb.FieldDescriptorProto,
	messagePb *descriptorpb.DescriptorProto,
) FieldDecoder {
	valueField := messagePb.GetField()[1]

	return FieldDecoder(fmt.Sprintf(
		"mapEntries \"%s\" %s",
		FieldJSONName(fieldPb),
		BasicFieldDecoder(nil, valueField),
	))
}

func MaybeType(t Type) Type {
	return Type(fmt.Sprintf("Maybe %s", t))
}

func MaybeEncoder(parentName *string, pb *descriptorpb.FieldDescriptorProto) FieldEncoder {
	return FieldEncoder(fmt.Sprintf(
		"optionalEncoder \"%s\" %s v.%s",
		FieldJSONName(pb),
		BasicFieldEncoder(parentName, pb),
		FieldName(pb.GetName()),
	))
}

func MaybeDecoder(parentName *string, pb *descriptorpb.FieldDescriptorProto) FieldDecoder {
	return FieldDecoder(fmt.Sprintf(
		"optional \"%s\" %s",
		FieldJSONName(pb),
		BasicFieldDecoder(parentName, pb),
	))
}

func ListType(t Type) Type {
	return Type(fmt.Sprintf("List %s", t))
}

func ListEncoder(parentName *string, pb *descriptorpb.FieldDescriptorProto) FieldEncoder {
	return FieldEncoder(fmt.Sprintf(
		"repeatedFieldEncoder \"%s\" %s v.%s",
		FieldJSONName(pb),
		BasicFieldEncoder(parentName, pb),
		FieldName(pb.GetName()),
	))
}

func ListDecoder(parentName *string, pb *descriptorpb.FieldDescriptorProto) FieldDecoder {
	return FieldDecoder(fmt.Sprintf(
		"repeated \"%s\" %s",
		FieldJSONName(pb),
		BasicFieldDecoder(parentName, pb),
	))
}

func OneOfType(preface []string, in string) Type {
	return NestedType(appendUnderscoreToReservedKeywords(stringextras.UpperCamelCase(in)), preface)
}

// TypeAliasTemplate - defines templates for self contained type aliases
func TypeAliasTemplate(t *template.Template) (*template.Template, error) {
	return t.Parse(`
{{- define "type-alias" -}}
type alias {{ .Name }} =
    { {{ range $i, $v := .Fields }}
        {{- if $i }}, {{ end }}{{ .Name }} : {{ .Type }}{{ if .Number }} -- {{ .Number }}{{ end }}
    {{ end }}}


{{ .Decoder }} : JD.Decoder {{ .Name }}
{{ .Decoder }} =
    JD.lazy <| \_ -> decode {{ .Name }}{{ range .Fields }}
        |> {{ .Decoder }}{{ end }}


{{ .Encoder }} : {{ .Name }} -> JE.Value
{{ .Encoder }} v =
    JE.object <| List.filterMap identity <|
        [{{ range $i, $v := .Fields }}
            {{- if $i }},{{ end }} ({{ .Encoder }})
        {{ end }}]
{{- end -}}
`)
}
