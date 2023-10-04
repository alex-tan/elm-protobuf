package parsepb

import (
	"strings"

	"github.com/thematthopkins/elm-protobuf/pkg/forwardextensions"
	"github.com/thematthopkins/elm-protobuf/pkg/generationparams"
	"google.golang.org/protobuf/proto"

	"github.com/thematthopkins/elm-protobuf/pkg/elm"
	"github.com/thematthopkins/elm-protobuf/pkg/stringextras"

	"google.golang.org/protobuf/types/descriptorpb"
)

type PbMessage struct {
	TypeAlias        elm.TypeAlias
	OneOfCustomTypes []elm.OneOfCustomType
	EnumCustomTypes  []elm.EnumCustomType
	NestedMessages   []PbMessage
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

func MessagesWithSingletons(messages []PbMessage) []PbMessage {
	result := []PbMessage{}

	for _, m := range messages {
		if m.TypeAlias.IsSingleton {
			result = append(result, m)
		}
	}
	return result
}

func MessagesWithIds(messages []PbMessage) []PbMessage {
	result := []PbMessage{}

	for _, m := range messages {
		has_id := false
		for _, f := range m.TypeAlias.Fields {
			if f.Name == "id" {
				has_id = true
			}
		}
		if has_id {
			result = append(result, m)
		}
	}
	return result
}

func EnumsToCustomTypes(preface []string, enumPbs []*descriptorpb.EnumDescriptorProto, p generationparams.Parameters) []elm.EnumCustomType {
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

func oneOfsToCustomTypes(preface []string, messagePb *descriptorpb.DescriptorProto, p generationparams.Parameters) []elm.OneOfCustomType {
	var result []elm.OneOfCustomType

	if isDeprecated(messagePb.Options) && p.RemoveDeprecated {
		return result
	}

	for oneofIndex, oneOfPb := range messagePb.GetOneofDecl() {
		syntheticField := syntheticFieldForOneOfIndex(messagePb, (int32)(oneofIndex))
		if syntheticField != nil {
			continue
		}

		variantPreface := append([]string{oneOfPb.GetName()}, preface...)
		var variants []elm.OneOfVariant
		for _, inField := range messagePb.GetField() {
			if isDeprecated(inField.Options) && p.RemoveDeprecated {
				continue
			}

			if inField.OneofIndex == nil || inField.GetOneofIndex() != int32(oneofIndex) {
				continue
			}

			variants = append(variants, elm.OneOfVariant{
				Name:     elm.NestedVariantName(inField.GetName(), variantPreface),
				JSONName: elm.OneOfVariantJSONName(inField),
				Type:     elm.BasicFieldType(nil, inField),
				Decoder:  elm.BasicFieldDecoder(nil, inField),
				Encoder:  elm.BasicFieldEncoder(nil, inField),
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

func syntheticFieldForOneOfIndex(messagePb *descriptorpb.DescriptorProto, oneofIndex int32) *descriptorpb.FieldDescriptorProto {
	for _, field := range messagePb.GetField() {
		if field.GetProto3Optional() && field.GetOneofIndex() == int32(oneofIndex) {
			return field
		}
	}
	return nil
}

func Messages(preface []string, messagePbs []*descriptorpb.DescriptorProto, p generationparams.Parameters) []PbMessage {
	var result []PbMessage

	for _, messagePb := range messagePbs {
		isSingleton := false
		if proto.HasExtension(messagePb.Options, forwardextensions.E_Singleton) {
			isSingleton = true
		}
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
			if nested != nil {
				newFields = append(newFields, elm.TypeAliasField{
					Name:           elm.FieldName(fieldPb.GetName()),
					IdTypeOverride: elm.GetIdTypeOverride(fieldPb),
					Type:           elm.MapType(nested),
					Number:         elm.ProtobufFieldNumber(fieldPb.GetNumber()),
					Encoder:        elm.MapEncoder(fieldPb, nested),
					Decoder:        elm.MapDecoder(fieldPb, nested),
				})
			} else if isOptional(fieldPb) {
				newFields = append(newFields, elm.TypeAliasField{
					Name:           elm.FieldName(fieldPb.GetName()),
					IdTypeOverride: elm.GetIdTypeOverride(fieldPb),
					Type:           elm.MaybeType(elm.BasicFieldType(messagePb.Name, fieldPb)),
					Number:         elm.ProtobufFieldNumber(fieldPb.GetNumber()),
					Encoder:        elm.MaybeEncoder(messagePb.Name, fieldPb),
					Decoder:        elm.MaybeDecoder(messagePb.Name, fieldPb),
				})
			} else if isRepeated(fieldPb) {
				newFields = append(newFields, elm.TypeAliasField{
					Name:           elm.FieldName(fieldPb.GetName()),
					IdTypeOverride: elm.GetIdTypeOverride(fieldPb),
					Type:           elm.ListType(elm.BasicFieldType(messagePb.Name, fieldPb)),
					Number:         elm.ProtobufFieldNumber(fieldPb.GetNumber()),
					Encoder:        elm.ListEncoder(messagePb.Name, fieldPb),
					Decoder:        elm.ListDecoder(messagePb.Name, fieldPb),
				})
			} else {
				newFields = append(newFields, elm.TypeAliasField{
					Name:           elm.FieldName(fieldPb.GetName()),
					IdTypeOverride: elm.GetIdTypeOverride(fieldPb),
					Type:           elm.BasicFieldType(messagePb.Name, fieldPb),
					Number:         elm.ProtobufFieldNumber(fieldPb.GetNumber()),
					Encoder:        elm.RequiredFieldEncoder(messagePb.Name, fieldPb),
					Decoder:        elm.RequiredFieldDecoder(messagePb.Name, fieldPb),
				})
			}
		}

		newPreface := append([]string{messagePb.GetName()}, preface...)
		for oneofIndex, oneOfPb := range messagePb.GetOneofDecl() {
			syntheticField := syntheticFieldForOneOfIndex(messagePb, (int32)(oneofIndex))
			if syntheticField != nil {
				newFields = append(newFields, elm.TypeAliasField{
					Name:    elm.FieldName(syntheticField.GetName()),
					Type:    elm.MaybeType(elm.BasicFieldType(nil, syntheticField)),
					Encoder: elm.MaybeEncoder(nil, syntheticField),
					Decoder: elm.MaybeDecoder(nil, syntheticField),
				})
			} else {
				newFields = append(newFields, elm.TypeAliasField{
					Name:    elm.FieldName(oneOfPb.GetName()),
					Type:    elm.OneOfType(newPreface, oneOfPb.GetName()),
					Encoder: elm.OneOfEncoder(newPreface, oneOfPb),
					Decoder: elm.OneOfDecoder(newPreface, oneOfPb),
				})
			}
		}

		name := elm.NestedType(messagePb.GetName(), preface)
		result = append(result, PbMessage{
			TypeAlias: elm.TypeAlias{
				IsSingleton: isSingleton,
				Name:        name,
				LowerName:   stringextras.FirstLower((string)(name)),
				Decoder:     elm.DecoderName(name),
				Encoder:     elm.EncoderName(name),
				Fields:      newFields,
			},
			OneOfCustomTypes: oneOfsToCustomTypes(newPreface, messagePb, p),
			EnumCustomTypes:  EnumsToCustomTypes(newPreface, messagePb.GetEnumType(), p),
			NestedMessages:   Messages(newPreface, messagePb.GetNestedType(), p),
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
