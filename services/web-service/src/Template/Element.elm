module Template.Element exposing
    ( Element
    , ElementType(..)
    , Form
    , decoder
    , encode
    , initForm
    , typeFromString
    , typeOptions
    , typeToString
    )

import Json.Decode as Decode exposing (Decoder)
import Json.Decode.Pipeline as DecodePipeline
import Json.Encode as Encode
import Template.Font as Font exposing (Font)


type ElementId
    = ElementId Int


type alias Element =
    { id : ElementId
    , x : Int
    , y : Int
    , width : Int
    , height : Int
    , type_ : ElementType
    }


type alias Form =
    { x : Int
    , y : Int
    , type_ : ElementType
    }


type ElementType
    = Text { font : Font, fontSize : Int, valueFrom : String }
    | Rect { width : Int, height : Int }


initForm : Form
initForm =
    { x = 0
    , y = 0
    , type_ = Rect { width = 5, height = 5 }
    }


typeOptions : List ( String, String )
typeOptions =
    [ Rect { width = 0, height = 0}, Text { font = Font.default, fontSize = 0, valueFrom = ""} ]
        |> List.map (\t -> ( typeLabel t, typeToString t ))


typeToString : ElementType -> String
typeToString type_ =
    case type_ of
        Rect _ ->
            "rect"

        Text _ ->
            "text"


typeLabel : ElementType -> String
typeLabel type_ =
    case type_ of
        Rect _ ->
            "Rect"

        Text _ ->
            "Text"


typeFromString : String -> Maybe ElementType
typeFromString string =
    case string of
        "rect" ->
            Just (Rect { width = 0, height = 0 })

        "text" ->
            Just (Text { font = Font.default, fontSize = 0, valueFrom = "" })

        _ ->
            Nothing



-- Decode / Encode


decoder : Decoder Element
decoder =
    Decode.succeed Element
        |> DecodePipeline.required "id" idDecoder
        |> DecodePipeline.required "x" Decode.int
        |> DecodePipeline.required "y" Decode.int
        |> DecodePipeline.required "width" Decode.int
        |> DecodePipeline.required "height" Decode.int
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
                        Decode.map2 (\width height -> Rect { width = width, height = height })
                            (Decode.field "width" Decode.int)
                            (Decode.field "height" Decode.int)

                    "text" ->
                        Decode.map3 (\font fontSize valueFrom -> Text { font = font, fontSize = fontSize, valueFrom = valueFrom })
                            (Decode.field "font" Font.decoder)
                            (Decode.field "font_size" Decode.int)
                            (Decode.field "value_from" Decode.string)

                    _ ->
                        Decode.fail ("Unknown element type: " ++ string)
            )


encode : Form -> Encode.Value
encode element =
    Encode.object
        ([ ( "x", Encode.int element.x )
         , ( "y", Encode.int element.y )
         ]
            ++ encodeType element.type_
        )


encodeType : ElementType -> List ( String, Encode.Value )
encodeType type_ =
    case type_ of
        Rect { width, height } ->
            [ ( "type", Encode.string (typeToString type_) )
            , ( "width", Encode.int width )
            , ( "height", Encode.int height )
            ]

        Text { font, fontSize, valueFrom } ->
            [ ( "type", Encode.string (typeToString type_) )
            , ( "font", Font.encode font )
            , ( "font_size", Encode.int fontSize )
            , ( "value_from", Encode.string valueFrom )
            ]
