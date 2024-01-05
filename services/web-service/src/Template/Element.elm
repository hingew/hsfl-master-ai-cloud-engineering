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


typeOptions : List ( String, String )
typeOptions =
    [ ( "Rect", "react" )
    , ( "Text", "text" )
    ]


typeToString : ElementType -> String
typeToString type_ =
    case type_ of
        Rect ->
            "rect"

        Text _ _ ->
            "text"


typeFromString : String -> Maybe ElementType
typeFromString string =
    case string of
        "rect" ->
            Just Rect

        "text" ->
            Just (Text "" 0)

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
                            (Decode.field "font_size" Decode.int)

                    _ ->
                        Decode.fail ("Unknown element type: " ++ string)
            )


encode : Form -> Encode.Value
encode element =
    Encode.object
        ([ ( "x", Encode.int element.x )
         , ( "y", Encode.int element.y )
         , ( "width", Encode.int element.width )
         , ( "height", Encode.int element.height )
         , ( "value_from", Encode.string element.valueFrom )
         ]
            ++ encodeType element.type_
        )


encodeType : ElementType -> List ( String, Encode.Value )
encodeType type_ =
    case type_ of
        Rect ->
            [ ( "type", Encode.string (typeToString type_) )
            ]

        Text font fontSize ->
            [ ( "type", Encode.string (typeToString type_) )
            , ( "font", Encode.string font )
            , ( "font_size", Encode.int fontSize )
            ]
