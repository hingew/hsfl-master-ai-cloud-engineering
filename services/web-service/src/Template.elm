module Template exposing
    ( CreateResponse
    , Template
    , TemplateId
    , create
    , decoder
    , delete
    , fetchAll
    , id
    , toId
    )

import Api
import Http
import Iso8601
import Json.Decode as Decode exposing (Decoder)
import Json.Decode.Pipeline as DecodePipeline
import Json.Encode as Encode
import RemoteData exposing (WebData)
import Session exposing (Token)
import Template.Element as Element exposing (Element)
import Time


type TemplateId
    = TemplateId Int


type alias Template =
    { id : TemplateId
    , createdAt : Time.Posix
    , updatedAt : Time.Posix
    , name : String
    , elements : List Element
    }



-- Utils


toId : Int -> TemplateId
toId =
    TemplateId


id : TemplateId -> Int
id (TemplateId value) =
    value



-- Decoder / Encode


decoder : Decoder Template
decoder =
    Decode.succeed Template
        |> DecodePipeline.required "id" idDecoder
        |> DecodePipeline.required "created_at" Iso8601.decoder
        |> DecodePipeline.required "updated_at" Iso8601.decoder
        |> DecodePipeline.required "name" Decode.string
        |> DecodePipeline.required "elements" (Decode.list Element.decoder)


createDecoder : Decoder CreateResponse
createDecoder =
    Decode.map CreateResponse
        (Decode.field "id" idDecoder)


idDecoder : Decoder TemplateId
idDecoder =
    Decode.int |> Decode.map TemplateId


encode : TemplateRequest -> Encode.Value
encode template =
    Encode.object
        [ ( "name", Encode.string template.name )
        , ( "elements", Encode.list Element.encode template.elements )
        ]



-- HTTP requests


fetchAll : Token -> (WebData (List Template) -> msg) -> Cmd msg
fetchAll token msg =
    Api.get
        { url = path
        , expect = Http.expectJson (RemoteData.fromResult >> msg) (Decode.list decoder)
        , token = token
        }


type alias TemplateRequest =
    { name : String
    , elements : List Element.Form
    }


type alias CreateResponse =
    { id : TemplateId }


create : Token -> TemplateRequest -> (WebData CreateResponse -> msg) -> Cmd msg
create token template msg =
    Api.post
        { url = path
        , expect = Http.expectJson (RemoteData.fromResult >> msg) createDecoder
        , body = Http.jsonBody (encode template)
        , token = token
        }


delete : Token -> TemplateId -> (WebData () -> msg) -> Cmd msg
delete token (TemplateId templateId) msg =
    Api.delete
        { url = path ++ "/" ++ String.fromInt templateId
        , expect = Http.expectWhatever (RemoteData.fromResult >> msg)
        , token = token
        }


path : String
path =
    "/api/templates"
