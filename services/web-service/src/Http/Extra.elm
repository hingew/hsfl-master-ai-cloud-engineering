module Http.Extra exposing (errToString)

import Http


errToString : Http.Error -> String
errToString err =
    case err of
        Http.BadUrl string ->
            "The url " ++ string ++ " is invalid!"

        Http.Timeout ->
            "The request timed out!"

        Http.NetworkError ->
            "There was an network error..."

        Http.BadStatus status ->
            "BadStatus:  " ++ String.fromInt status

        Http.BadBody string ->
            "BadBody: " ++ string
