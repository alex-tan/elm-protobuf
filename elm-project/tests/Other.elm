module Other exposing (..)

-- DO NOT EDIT
-- AUTOGENERATED BY THE ELM PROTOCOL BUFFER COMPILER
-- https://github.com/tiziano88/elm-protobuf
-- source file: other.proto

import Protobuf exposing (..)

import Json.Decode as JD
import Json.Encode as JE


uselessDeclarationToPreventErrorDueToEmptyOutputFile = 42

requiredWithoutDefault : String -> JD.Decoder a -> JD.Decoder (a -> b) -> JD.Decoder b
requiredWithoutDefault name decoder d =
    field (JD.field name decoder) d

requiredFieldEncoderWithoutDefault : String -> (a -> JE.Value) -> a -> Maybe ( String, JE.Value )
requiredFieldEncoderWithoutDefault name encoder v =
    Just ( name, encoder v )


type alias Other =
    { stringField : String -- 1
    }


otherDecoder : JD.Decoder Other
otherDecoder =
    JD.lazy <| \_ -> decode Other
        |> required "stringField" JD.string ""


otherEncoder : Other -> JE.Value
otherEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "stringField" JE.string "" v.stringField)
        ]
