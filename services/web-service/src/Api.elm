module Api exposing
    ( delete
    , get
    , post
    , put
    )

import Http exposing (Body)
import Session exposing (Token)


get : { url : String, expect : Http.Expect msg, token : Token } -> Cmd msg
get { url, expect, token } =
    Http.request
        { method = "GET"
        , headers = [ Http.header "Authorization" ("Bearer " ++ Session.tokenValue token) ]
        , body = Http.emptyBody
        , url = url
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }


post : { url : String, expect : Http.Expect msg, body : Body, token : Token } -> Cmd msg
post { url, expect, body, token } =
    Http.request
        { method = "POST"
        , headers = [ Http.header "Authorization" ("Bearer " ++ Session.tokenValue token) ]
        , body = body
        , url = url
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }


put : { url : String, expect : Http.Expect msg, token : Token, body : Body } -> Cmd msg
put { url, expect, token, body } =
    Http.request
        { method = "PUT"
        , headers = [ Http.header "Authorization" ("Bearer " ++ Session.tokenValue token) ]
        , body = body
        , url = url
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }


delete : { url : String, expect : Http.Expect msg, token : Token } -> Cmd msg
delete { url, expect, token } =
    Http.request
        { method = "DELETE"
        , headers = [ Http.header "Authorization" ("Bearer " ++ Session.tokenValue token) ]
        , body = Http.emptyBody
        , url = url
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }
