module DomainNew.Graph exposing (..)

-- DO NOT EDIT
-- AUTOGENERATED BY THE ELM PROTOCOL BUFFER COMPILER
-- https://github.com/tiziano88/elm-protobuf
-- source file: forward.proto


import ForwardNew.Interface.Cache as Cache
import ForwardNew.Lookup as Lookup exposing (Lookup)
import Ids
import Json.Decode as Decode
import LocalExtra.Lookup as LookupExtra
import Pb

allDecoders : List Lookup.DecoderConfig
allDecoders = [
	 Lookup.toDecoderConfig myEntity 
	, Lookup.toDecoderConfig myChildEntity 
	, Lookup.toDecoderConfig unreferencedEntity 
	, Lookup.toDecoderConfig selfReferencing 
	, Lookup.toDecoderConfig overrideName 
]


myEntity : Lookup Ids.MyEntity Pb.MyEntity
myEntity =
    Lookup.defineNode
        { entrypoint = "myEntity"
        , parameters = LookupExtra.idParam (\(Ids.MyEntity id) -> id)
        , decoder = Pb.myEntity
        , cacheKey = Cache.myEntity
        }

myChildEntity : Lookup Ids.MyChildEntity Pb.MyChildEntity
myChildEntity =
    Lookup.defineNode
        { entrypoint = "myChildEntity"
        , parameters = LookupExtra.idParam (\(Ids.MyChildEntity id) -> id)
        , decoder = Pb.myChildEntity
        , cacheKey = Cache.myChildEntity
        }

unreferencedEntity : Lookup Ids.UnreferencedEntity Pb.UnreferencedEntity
unreferencedEntity =
    Lookup.defineNode
        { entrypoint = "unreferencedEntity"
        , parameters = LookupExtra.idParam (\(Ids.UnreferencedEntity id) -> id)
        , decoder = Pb.unreferencedEntity
        , cacheKey = Cache.unreferencedEntity
        }

selfReferencing : Lookup Ids.SelfReferencing Pb.SelfReferencing
selfReferencing =
    Lookup.defineNode
        { entrypoint = "selfReferencing"
        , parameters = LookupExtra.idParam (\(Ids.SelfReferencing id) -> id)
        , decoder = Pb.selfReferencing
        , cacheKey = Cache.selfReferencing
        }

overrideName : Lookup Ids.OverrideName Pb.OverrideName
overrideName =
    Lookup.defineNode
        { entrypoint = "overrideName"
        , parameters = LookupExtra.idParam (\(Ids.OverrideName id) -> id)
        , decoder = Pb.overrideName
        , cacheKey = Cache.overrideName
        }

