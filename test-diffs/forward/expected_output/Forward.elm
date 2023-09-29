module Forward exposing (..)

-- DO NOT EDIT
-- AUTOGENERATED BY THE ELM PROTOCOL BUFFER COMPILER
-- https://github.com/tiziano88/elm-protobuf
-- source file: forward.proto

import Protobuf exposing (..)

import Json.Decode as JD
import Json.Encode as JE


uselessDeclarationToPreventErrorDueToEmptyOutputFile = 42


type alias MySingleton =
    { myChildEntityIds : List Ids.MyChildEntity -- 2
    }


mySingletonDecoder : JD.Decoder MySingleton
mySingletonDecoder =
    JD.lazy <| \_ -> decode MySingleton
        |> repeated "myChildEntityIds" (JD.string |> JD.map Ids.MyChildEntity)


mySingletonEncoder : MySingleton -> JE.Value
mySingletonEncoder v =
    JE.object <| List.filterMap identity <|
        [ (repeatedFieldEncoder "myChildEntityIds" (\(Ids.MyChildEntity id) -> JE.string id) v.myChildEntityIds)
        ]


type alias MyEntity =
    { id : Ids.MyEntity -- 1
    , myChildEntityIds : List Ids.MyChildEntity -- 2
    }


myEntityDecoder : JD.Decoder MyEntity
myEntityDecoder =
    JD.lazy <| \_ -> decode MyEntity
        |> required "id" (JD.string |> JD.map Ids.MyEntity) (Ids.MyEntity "")
        |> repeated "myChildEntityIds" (JD.string |> JD.map Ids.MyChildEntity)


myEntityEncoder : MyEntity -> JE.Value
myEntityEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "id" (\(Ids.MyEntity id) -> JE.string id) (Ids.MyEntity "") v.id)
        , (repeatedFieldEncoder "myChildEntityIds" (\(Ids.MyChildEntity id) -> JE.string id) v.myChildEntityIds)
        ]


type alias MyChildEntity =
    { id : Ids.MyChildEntity -- 1
    }


myChildEntityDecoder : JD.Decoder MyChildEntity
myChildEntityDecoder =
    JD.lazy <| \_ -> decode MyChildEntity
        |> required "id" (JD.string |> JD.map Ids.MyChildEntity) (Ids.MyChildEntity "")


myChildEntityEncoder : MyChildEntity -> JE.Value
myChildEntityEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "id" (\(Ids.MyChildEntity id) -> JE.string id) (Ids.MyChildEntity "") v.id)
        ]


type alias UnreferencedEntity =
    { id : Ids.UnreferencedEntity -- 1
    }


unreferencedEntityDecoder : JD.Decoder UnreferencedEntity
unreferencedEntityDecoder =
    JD.lazy <| \_ -> decode UnreferencedEntity
        |> required "id" (JD.string |> JD.map Ids.UnreferencedEntity) (Ids.UnreferencedEntity "")


unreferencedEntityEncoder : UnreferencedEntity -> JE.Value
unreferencedEntityEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "id" (\(Ids.UnreferencedEntity id) -> JE.string id) (Ids.UnreferencedEntity "") v.id)
        ]


type alias SelfReferencing =
    { id : Ids.SelfReferencing -- 1
    , selfReferencingId : Ids.SelfReferencing -- 2
    }


selfReferencingDecoder : JD.Decoder SelfReferencing
selfReferencingDecoder =
    JD.lazy <| \_ -> decode SelfReferencing
        |> required "id" (JD.string |> JD.map Ids.SelfReferencing) (Ids.SelfReferencing "")
        |> required "selfReferencingId" (JD.string |> JD.map Ids.SelfReferencing) (Ids.SelfReferencing "")


selfReferencingEncoder : SelfReferencing -> JE.Value
selfReferencingEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "id" (\(Ids.SelfReferencing id) -> JE.string id) (Ids.SelfReferencing "") v.id)
        , (requiredFieldEncoder "selfReferencingId" (\(Ids.SelfReferencing id) -> JE.string id) (Ids.SelfReferencing "") v.selfReferencingId)
        ]


type alias OverrideName =
    { id : Ids.OverrideName -- 1
    , referenceByOtherName : Ids.MyEntity -- 3
    , manyReferencesByOtherName : List Ids.MyEntity -- 4
    , optionalReferenceByOtherName : Maybe Ids.MyEntity
    }


overrideNameDecoder : JD.Decoder OverrideName
overrideNameDecoder =
    JD.lazy <| \_ -> decode OverrideName
        |> required "id" (JD.string |> JD.map Ids.OverrideName) (Ids.OverrideName "")
        |> required "referenceByOtherName" (JD.string |> JD.map Ids.MyEntity) (Ids.MyEntity "")
        |> repeated "manyReferencesByOtherName" (JD.string |> JD.map Ids.MyEntity)
        |> optional "optionalReferenceByOtherName" (JD.string |> JD.map Ids.MyEntity)


overrideNameEncoder : OverrideName -> JE.Value
overrideNameEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "id" (\(Ids.OverrideName id) -> JE.string id) (Ids.OverrideName "") v.id)
        , (requiredFieldEncoder "referenceByOtherName" (\(Ids.MyEntity id) -> JE.string id) (Ids.MyEntity "") v.referenceByOtherName)
        , (repeatedFieldEncoder "manyReferencesByOtherName" (\(Ids.MyEntity id) -> JE.string id) v.manyReferencesByOtherName)
        , (optionalEncoder "optionalReferenceByOtherName" (\(Ids.MyEntity id) -> JE.string id) v.optionalReferenceByOtherName)
        ]


type alias MessageWithoutId =
    { noIdHere : String -- 1
    }


messageWithoutIdDecoder : JD.Decoder MessageWithoutId
messageWithoutIdDecoder =
    JD.lazy <| \_ -> decode MessageWithoutId
        |> required "noIdHere" JD.string ""


messageWithoutIdEncoder : MessageWithoutId -> JE.Value
messageWithoutIdEncoder v =
    JE.object <| List.filterMap identity <|
        [ (requiredFieldEncoder "noIdHere" JE.string "" v.noIdHere)
        ]
