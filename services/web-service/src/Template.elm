module Template exposing (Template, decoder, fetchAll)

import Http
import Iso8601
import Json.Decode as Decode exposing (Decoder)
import Json.Decode.Pipeline as DecodePipeline
import RemoteData exposing (WebData)
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


decoder : Decoder Template
decoder =
    Decode.succeed Template
        |> DecodePipeline.required "id" idDecoder
        |> DecodePipeline.required "created_at" Iso8601.decoder
        |> DecodePipeline.required "updated_at" Iso8601.decoder
        |> DecodePipeline.required "name" Decode.string
        |> DecodePipeline.required "elements" (Decode.list Element.decoder)


idDecoder : Decoder TemplateId
idDecoder =
    Decode.int |> Decode.map TemplateId


fetchAll : (WebData (List Template) -> msg) -> Cmd msg
fetchAll msg =
    Http.get
        { url = path
        , expect = Http.expectJson (RemoteData.fromResult >> msg) (Decode.list decoder)
        }


path : String
path =
    "/api/templates"
