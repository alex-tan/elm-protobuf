module Simple exposing (..)

-- DO NOT EDIT
-- AUTOGENERATED BY THE ELM PROTOCOL BUFFER COMPILER
-- https://github.com/tiziano88/elm-protobuf
-- source file: simple.proto

import Protobuf exposing (..)
import Dict exposing (Dict)
import Json.Decode as JD
import Json.Encode as JE
import Dir.Other_dir exposing (..)

import Other exposing (..)



uselessDeclarationToPreventErrorDueToEmptyOutputFile = 42

requiredWithoutDefault : String -> JD.Decoder a -> JD.Decoder (a -> b) -> JD.Decoder b
requiredWithoutDefault name decoder d =
    field (JD.field name decoder) d

requiredFieldEncoderWithoutDefault : String -> (a -> JE.Value) -> a -> Maybe ( String, JE.Value )
requiredFieldEncoderWithoutDefault name encoder v =
    Just ( name, encoder v )


type Colour
    = ColourUnspecified -- 0
    | Red -- 1
    | Green -- 2
    | Blue -- 3


colourDecoder : JD.Decoder Colour
colourDecoder =
    JD.map (Maybe.withDefault colourDefault << colourFromString) JD.string


colourDefault : Colour
colourDefault = ColourUnspecified


colourToString : Colour -> String
colourToString v =
    case v of
        ColourUnspecified ->
            "COLOUR_UNSPECIFIED"

        Red ->
            "RED"

        Green ->
            "GREEN"

        Blue ->
            "BLUE"


allColour : List Colour
allColour =
  [ ColourUnspecified
  , Red
  , Green
  , Blue
  ]

colourDict : Dict String Colour
colourDict =
    Dict.fromList <|
        List.map
            (\v -> ( colourToString v, v ))
            allColour

colourFromString : String -> Maybe Colour
colourFromString s =
    Dict.get s colourDict

colourEncoder : Colour -> JE.Value
colourEncoder =
    JE.string << colourToString


type alias Empty =
    { }


emptyDecoder : JD.Decoder Empty
emptyDecoder =
    JD.lazy <| \_ -> decode Empty


emptyEncoder : Empty -> JE.Value
emptyEncoder v =
    JE.object <| List.filterMap identity <|
        []


type alias Simple =
    { int32Field : Int -- 1
    }


simpleDecoder : JD.Decoder Simple
simpleDecoder =
    JD.lazy <| \_ -> decode Simple
        |> required "int32Field" intDecoder 0


simpleEncoder : Simple -> JE.Value
simpleEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "int32Field" JE.int 0 v.int32Field)
        ]


type alias Foo =
    { s : Simple -- 1
    , ss : List Simple -- 2
    , colour : Colour -- 3
    , colours : List Colour -- 4
    , singleIntField : Int -- 5
    , repeatedIntField : List Int -- 6
    , bytesField : Bytes -- 9
    , oo : Foo_Oo
    , stringValueField : Maybe String
    , otherField : Maybe Other
    , otherDirField : Maybe OtherDir
    , timestampField : Maybe Timestamp
    }


fooDecoder : JD.Decoder Foo
fooDecoder =
    JD.lazy <| \_ -> decode Foo
        |> requiredWithoutDefault "s" simpleDecoder
        |> repeated "ss" simpleDecoder
        |> required "colour" colourDecoder colourDefault
        |> repeated "colours" colourDecoder
        |> required "singleIntField" intDecoder 0
        |> repeated "repeatedIntField" intDecoder
        |> required "bytesField" bytesFieldDecoder []
        |> field foo_OoDecoder
        |> optional "stringValueField" stringValueDecoder
        |> optional "otherField" otherDecoder
        |> optional "otherDirField" otherDirDecoder
        |> optional "timestampField" timestampDecoder


fooEncoder : Foo -> JE.Value
fooEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoderWithoutDefault "s" simpleEncoder v.s)
        , (repeatedFieldEncoder "ss" simpleEncoder v.ss)
        , (requiredFieldEncoder "colour" colourEncoder colourDefault v.colour)
        , (repeatedFieldEncoder "colours" colourEncoder v.colours)
        , (requiredFieldEncoder "singleIntField" JE.int 0 v.singleIntField)
        , (repeatedFieldEncoder "repeatedIntField" JE.int v.repeatedIntField)
        , (requiredFieldEncoder "bytesField" bytesFieldEncoder [] v.bytesField)
        , (foo_OoEncoder v.oo)
        , (optionalEncoder "stringValueField" stringValueEncoder v.stringValueField)
        , (optionalEncoder "otherField" otherEncoder v.otherField)
        , (optionalEncoder "otherDirField" otherDirEncoder v.otherDirField)
        , (optionalEncoder "timestampField" timestampEncoder v.timestampField)
        ]


type Foo_Oo
    = 
     Foo_Oo_Oo1 Int
    | Foo_Oo_Oo2 Bool


foo_OoDecoder : JD.Decoder Foo_Oo
foo_OoDecoder =
    JD.lazy <| \_ -> JD.oneOf
        [ JD.map Foo_Oo_Oo1 (JD.field "oo1" intDecoder)
        , JD.map Foo_Oo_Oo2 (JD.field "oo2" JD.bool)
        , JD.fail "Foo_Oo_Unspecified"
        ]


foo_OoEncoder : Foo_Oo -> Maybe ( String, JE.Value )
foo_OoEncoder v =
    case v of

        Foo_Oo_Oo1 x ->
            Just ( "oo1", JE.int x )

        Foo_Oo_Oo2 x ->
            Just ( "oo2", JE.bool x )
