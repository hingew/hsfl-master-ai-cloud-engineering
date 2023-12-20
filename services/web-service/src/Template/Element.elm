module Template.Element exposing
    ( Element
    , Form
    , decoder
    , initForm
    )

import Json.Decode as Decode exposing (Decoder)
import Json.Decode.Pipeline as DecodePipeline


type ElementId
    = ElementId Int


type alias Element =
    { id : ElementId
    , x : Int
    , y : Int
    , width : Int
    , height : Int
    , valueFrom : String
    , type_ : ElementType
    }


type alias Form =
    { x : Int
    , y : Int
    , width : Int
    , height : Int
    , valueFrom : String
    , type_ : ElementType
    }


type ElementType
    = Text String Int
    | Rect


initForm : Form
initForm =
    { x = 0
    , y = 0
    , width = 5
    , height = 5
    , valueFrom = "value"
    , type_ = Rect
    }



-- Decode / Encode


decoder : Decoder Element
decoder =
    Decode.succeed Element
        |> DecodePipeline.required "id" idDecoder
        |> DecodePipeline.required "x" Decode.int
        |> DecodePipeline.required "y" Decode.int
        |> DecodePipeline.required "width" Decode.int
        |> DecodePipeline.required "height" Decode.int
        |> DecodePipeline.required "value_from" Decode.string
        |> DecodePipeline.custom typeDecoder


idDecoder : Decoder ElementId
idDecoder =
    Decode.int |> Decode.map ElementId


typeDecoder : Decoder ElementType
typeDecoder =
    Decode.field "type" Decode.string
        |> Decode.andThen
            (\string ->
                case string of
                    "rect" ->
                        Decode.succeed Rect

                    "text" ->
                        Decode.map2 Text
                            (Decode.field "font" Decode.string)
                            (Decode.field "fontSize" Decode.int)

                    _ ->
                        Decode.fail ("Unknown element type: " ++ string)
            )
