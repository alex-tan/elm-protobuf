package forwardids

import (
	"bytes"
	"github.com/thematthopkins/elm-protobuf/pkg/parsepb"
	"github.com/thematthopkins/elm-protobuf/pkg/stringextras"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"sort"
	"strings"
	"text/template"
)

func Generate(inFile *descriptorpb.FileDescriptorProto, messages []parsepb.PbMessage) (*pluginpb.CodeGeneratorResponse_File, error) {
	t := template.New("t")

	t, err := t.Parse(`module Ids exposing (..)

-- DO NOT EDIT
-- AUTOGENERATED BY THE ELM PROTOCOL BUFFER COMPILER
-- https://github.com/tiziano88/elm-protobuf
-- source file: {{ .SourceFile }}

{{ range $index, $element := .IdTypes}}
type {{$element.UpperName}}
    = {{$element.UpperName}} String

{{$element.LowerName}}ToString : {{$element.UpperName}} -> String
{{$element.LowerName}}ToString ({{$element.UpperName}} v) = v

{{end}}

`)
	if err != nil {
		return nil, err
	}

	idTypes := map[string]struct{}{}
	IdTypesForMessages(messages, &idTypes)

	idTypesArr := []string{}
	for s := range idTypes {
		idTypesArr = append(idTypesArr, s)
	}

	sort.Strings(idTypesArr)

	idTypeStructs := []IdType{}

	for _, s := range idTypesArr {
		idTypeStructs = append(idTypeStructs, IdType{
			UpperName: s,
			LowerName: stringextras.LowerCamelCase(s),
		})
	}

	buff := &bytes.Buffer{}
	if err = t.Execute(buff, struct {
		SourceFile string
		IdTypes    []IdType
	}{
		SourceFile: inFile.GetName(),
		IdTypes:    idTypeStructs,
	}); err != nil {
		return nil, err
	}

	fileName := "Ids.elm"
	result := buff.String()
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &fileName,
		Content: &result,
	}, nil
}

type IdType struct {
	UpperName string
	LowerName string
}

func IdTypesForMessages(messages []parsepb.PbMessage, result *map[string]struct{}) {
	for _, m := range messages {
		for _, f := range m.TypeAlias.Fields {
			if strings.HasPrefix((string)(f.Type), "Ids.") {
				(*result)[strings.TrimPrefix((string)(f.Type), "Ids.")] = struct{}{}
			}
		}

		IdTypesForMessages(m.NestedMessages, result)
	}
}
