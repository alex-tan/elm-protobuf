package forwardgraph

import (
	"bytes"
	"fmt"
	"github.com/thematthopkins/elm-protobuf/pkg/elmpb"
	"github.com/thematthopkins/elm-protobuf/pkg/parsepb"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"text/template"
)

func Generate(inFile *descriptorpb.FileDescriptorProto, messages []parsepb.PbMessage) (*pluginpb.CodeGeneratorResponse_File, error) {
	t := template.New("t")

	t, err := t.Parse(`module DomainNew.Graph exposing (..)

-- DO NOT EDIT
-- AUTOGENERATED BY THE ELM PROTOCOL BUFFER COMPILER
-- https://github.com/tiziano88/elm-protobuf
-- source file: {{ .SourceFile }}


import ForwardNew.Interface.Cache as Cache
import ForwardNew.Lookup as Lookup exposing (Lookup)
import Ids
import Json.Decode as Decode
import LocalExtra.Lookup as LookupExtra
import {{.MainPackage}}

allDecoders : List Lookup.DecoderConfig
allDecoders = [{{ range $index, $element := .Decoders}}
    {{if $index}},{{end}} Lookup.toDecoderConfig {{ $element.Entrypoint }} {{end}}
    ]

{{ range .Decoders}}
{{ .Entrypoint }} : Lookup {{.LookupIdType}} Pb.{{ .UpperName }}
{{ .Entrypoint }} =
    Lookup.defineNode
        { entrypoint = "{{ .Entrypoint }}"
        , parameters = {{ .Parameters }}
        , decoder = Pb.{{ .LowerName }}Decoder
        , cacheKey = Cache.{{ .LowerName }}
        }
{{end}}
`)
	if err != nil {
		return nil, err
	}

	graph_messages := append(parsepb.MessagesWithSingletons(messages), parsepb.MessagesWithIds(messages)...)

	decoders := GetDecoders(graph_messages)

	mainPackage := elmpb.PackageName(inFile)

	buff := &bytes.Buffer{}
	if err = t.Execute(buff, struct {
		SourceFile  string
		MainPackage string
		Decoders    []Decoder
	}{
		SourceFile:  inFile.GetName(),
		MainPackage: mainPackage,
		Decoders:    decoders,
	}); err != nil {
		return nil, err
	}

	fileName := "DomainNew/Graph.elm"
	result := buff.String()
	return &pluginpb.CodeGeneratorResponse_File{
		Name:    &fileName,
		Content: &result,
	}, nil
}

type Decoder struct {
	UpperName    string
	LowerName    string
	Entrypoint   string
	Parameters   string
	LookupIdType string
	Decoder      string
	CacheKey     string
}

func GetDecoders(messages []parsepb.PbMessage) []Decoder {
	result := []Decoder{}

	for _, m := range messages {
		parameters := "Lookup.noParameters"
		lookupIdType := "()"
		entrypoint := (string)(m.TypeAlias.LowerName) + "Singleton"
		if !m.TypeAlias.IsSingleton {
			idType := ""
			for _, f := range m.TypeAlias.Fields {
				if f.Name == "id" {
					idType = (string)(f.Type)
				}
			}
			parameters = fmt.Sprintf("LookupExtra.idParam (\\(%s id) -> id)", idType)
			lookupIdType = idType
			entrypoint = (string)(m.TypeAlias.LowerName)
		}
		result = append(result, Decoder{
			UpperName:    (string)(m.TypeAlias.Name),
			LowerName:    m.TypeAlias.LowerName,
			Entrypoint:   entrypoint,
			LookupIdType: lookupIdType,
			Parameters:   parameters,
			Decoder:      fmt.Sprintf("Pb.%s", m.TypeAlias.Decoder),
			CacheKey:     fmt.Sprintf("Cache.%s", m.TypeAlias.LowerName),
		})
	}

	return result
}
