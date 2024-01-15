module Template.Font exposing (Font, decoder, default, encode, fromString, toString, options)

import Json.Decode as Decode exposing (Decoder)
import Json.Encode as Encode


type Font
    = Arial


default : Font
default =
    Arial


toString : Font -> String
toString font =
    case font of
        Arial ->
            "Arial"


fromString : String -> Maybe Font
fromString string =
    case string of
        "Arial" ->
            Just Arial

        _ ->
            Nothing


decoder : Decoder Font
decoder =
    Decode.string
        |> Decode.andThen
            (\string ->
                case fromString string of
                    Just font ->
                        Decode.succeed font

                    Nothing ->
                        Decode.fail ("unknown font: " ++ string)
            )


encode : Font -> Encode.Value
encode font =
    font
        |> toString
        |> Encode.string


options : List ( String, String )
options =
    [ Arial ]
        |> List.map
            (\font ->
                let
                    label =
                        toString font
                in
                ( label, label )
            )
