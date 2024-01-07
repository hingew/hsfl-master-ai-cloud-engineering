module Template exposing
    ( CreateResponse
    , PrintForm
    , Template
    , TemplateId
    , create
    , decoder
    , delete
    , fetch
    , fetchAll
    , id
    , print
    , setFormValue
    , toForm
    , toId
    )

import Api
import Bytes exposing (Bytes)
import Dict exposing (Dict)
import Http
import Http.Extra
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


type alias PrintForm =
    { id : TemplateId, values : Dict String String }



-- Utils


toId : Int -> TemplateId
toId =
    TemplateId


id : TemplateId -> Int
id (TemplateId value) =
    value


toForm : Template -> PrintForm
toForm template =
    { id = template.id
    , values =
        template.elements
            |> List.filterMap toFormValue
            |> Dict.fromList
    }


toFormValue : Element -> Maybe ( String, String )
toFormValue element =
    case element.type_ of
        Element.Rect _ _ ->
            Nothing

        Element.Text _ _ ->
            Just ( element.valueFrom, "" )


setFormValue : String -> String -> PrintForm -> PrintForm
setFormValue key value form =
    { form | values = Dict.insert key value form.values }



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


encodePrintForm : PrintForm -> Encode.Value
encodePrintForm { values } =
    values
        |> Dict.toList
        |> List.map (\( k, v ) -> ( k, Encode.string v ))
        |> Encode.object



-- HTTP requests


fetchAll : Token -> (WebData (List Template) -> msg) -> Cmd msg
fetchAll token msg =
    Api.get
        { url = path
        , expect = Http.expectJson (RemoteData.fromResult >> msg) (Decode.list decoder)
        , token = token
        }


fetch : TemplateId -> Token -> (WebData Template -> msg) -> Cmd msg
fetch templateId token msg =
    Api.get
        { url = detailPath templateId
        , expect = Http.expectJson (RemoteData.fromResult >> msg) decoder
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


print : Token -> PrintForm -> (WebData Bytes -> msg) -> Cmd msg
print token form msg =
    Api.post
        { url = "/api/render/" ++ String.fromInt (id form.id)
        , expect = Http.expectBytesResponse (RemoteData.fromResult >> msg) (Http.Extra.resolveBytes Ok)
        , body = Http.jsonBody (encodePrintForm form)
        , token = token
        }


delete : Token -> TemplateId -> (WebData () -> msg) -> Cmd msg
delete token templateId msg =
    Api.delete
        { url = detailPath templateId
        , expect = Http.expectWhatever (RemoteData.fromResult >> msg)
        , token = token
        }


path : String
path =
    "/api/templates"


detailPath : TemplateId -> String
detailPath (TemplateId templateId) =
    path ++ "/" ++ String.fromInt templateId
