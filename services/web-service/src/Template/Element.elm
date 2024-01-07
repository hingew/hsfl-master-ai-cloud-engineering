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
    , valueFrom : String
    , type_ : ElementType
    }


type ElementType
    = Text String Int
    | Rect Int Int


initForm : Form
initForm =
    { x = 0
    , y = 0
    , valueFrom = "value"
    , type_ = Rect 5 5
    }


typeOptions : List ( String, String )
typeOptions =
    [ ( "Rect", "rect" )
    , ( "Text", "text" )
    ]


typeToString : ElementType -> String
typeToString type_ =
    case type_ of
        Rect _ _ ->
            "rect"

        Text _ _ ->
            "text"


typeFromString : String -> Maybe ElementType
typeFromString string =
    case string of
        "rect" ->
            Just (Rect 0 0)

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
                        Decode.map2 Rect
                            (Decode.field "width" Decode.int)
                            (Decode.field "height" Decode.int)

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
         , ( "value_from", Encode.string element.valueFrom )
         ]
            ++ encodeType element.type_
        )


encodeType : ElementType -> List ( String, Encode.Value )
encodeType type_ =
    case type_ of
        Rect width height ->
            [ ( "type", Encode.string (typeToString type_) )
            , ( "width", Encode.int width )
            , ( "height", Encode.int height )
            ]

        Text font fontSize ->
            [ ( "type", Encode.string (typeToString type_) )
            , ( "font", Encode.string font )
            , ( "font_size", Encode.int fontSize )
            ]
