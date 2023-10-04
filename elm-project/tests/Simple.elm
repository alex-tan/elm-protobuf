module Simple exposing (..)

-- DO NOT EDIT
-- AUTOGENERATED BY THE ELM PROTOCOL BUFFER COMPILER
-- https://github.com/tiziano88/elm-protobuf
-- source file: simple.proto

import Protobuf exposing (..)

import Json.Decode as JD
import Json.Encode as JE
import Dir.Other_dir exposing (..)

import Other exposing (..)



uselessDeclarationToPreventErrorDueToEmptyOutputFile = 42


type Colour
    = ColourUnspecified -- 0
    | Red -- 1
    | Green -- 2
    | Blue -- 3


colourDecoder : JD.Decoder Colour
colourDecoder =
    let
        lookup s =
            case s of
                "COLOUR_UNSPECIFIED" ->
                    ColourUnspecified

                "RED" ->
                    Red

                "GREEN" ->
                    Green

                "BLUE" ->
                    Blue

                _ ->
                    ColourUnspecified
    in
        JD.map lookup JD.string


colourDefault : Colour
colourDefault = ColourUnspecified


colourEncoder : Colour -> JE.Value
colourEncoder v =
    let
        lookup s =
            case s of
                ColourUnspecified ->
                    "COLOUR_UNSPECIFIED"

                Red ->
                    "RED"

                Green ->
                    "GREEN"

                Blue ->
                    "BLUE"

    in
        JE.string <| lookup v


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
    { s : Maybe Simple -- 1
    , ss : List Simple -- 2
    , colour : Colour -- 3
    , colours : List Colour -- 4
    , singleIntField : Int -- 5
    , repeatedIntField : List Int -- 6
    , bytesField : Bytes -- 9
    , stringValueField : Maybe String -- 10
    , otherField : Maybe Other -- 11
    , otherDirField : Maybe OtherDir -- 12
    , timestampField : Maybe Timestamp -- 13
    , oo : Foo_Oo
    }


fooDecoder : JD.Decoder Foo
fooDecoder =
    JD.lazy <| \_ -> decode Foo
        |> optional "s" simpleDecoder
        |> repeated "ss" simpleDecoder
        |> required "colour" colourDecoder colourDefault
        |> repeated "colours" colourDecoder
        |> required "singleIntField" intDecoder 0
        |> repeated "repeatedIntField" intDecoder
        |> required "bytesField" bytesFieldDecoder []
        |> optional "stringValueField" stringValueDecoder
        |> optional "otherField" otherDecoder
        |> optional "otherDirField" otherDirDecoder
        |> optional "timestampField" timestampDecoder
        |> field foo_OoDecoder


fooEncoder : Foo -> JE.Value
fooEncoder v =
    JE.object <| List.filterMap identity <|
        [ (optionalEncoder "s" simpleEncoder v.s)
        , (repeatedFieldEncoder "ss" simpleEncoder v.ss)
        , (requiredFieldEncoder "colour" colourEncoder colourDefault v.colour)
        , (repeatedFieldEncoder "colours" colourEncoder v.colours)
        , (requiredFieldEncoder "singleIntField" JE.int 0 v.singleIntField)
        , (repeatedFieldEncoder "repeatedIntField" JE.int v.repeatedIntField)
        , (requiredFieldEncoder "bytesField" bytesFieldEncoder [] v.bytesField)
        , (optionalEncoder "stringValueField" stringValueEncoder v.stringValueField)
        , (optionalEncoder "otherField" otherEncoder v.otherField)
        , (optionalEncoder "otherDirField" otherDirEncoder v.otherDirField)
        , (optionalEncoder "timestampField" timestampEncoder v.timestampField)
        , (foo_OoEncoder v.oo)
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
