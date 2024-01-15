module Http.Extra exposing (errToString, resolveBytes)

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

resolveBytes : (body -> Result String a) -> Http.Response body -> Result Http.Error a
resolveBytes toResult response =
    case response of
        Http.BadUrl_ url_ ->
            Err (Http.BadUrl url_)

        Http.Timeout_ ->
            Err Http.Timeout

        Http.NetworkError_ ->
            Err Http.NetworkError

        Http.BadStatus_ metadata _ ->
            Err (Http.BadStatus metadata.statusCode)

        Http.GoodStatus_ _ body ->
            Result.mapError Http.BadBody (toResult body)

