module Api exposing
    ( get
    , post
    , put
    )

import Http exposing (Body)
import Session exposing (Token)


get : { url : String, expect : Http.Expect msg , token : Token} -> Cmd msg
get { url, expect, token } =
    Http.request
        { method = "GET"
        , headers = [ Http.header "Authorization" ("Bearer " ++ (Session.tokenValue token)) ]
        , body = Http.emptyBody
        , url = url
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }


post : { url : String, expect : Http.Expect msg, body : Body } -> Cmd msg
post { url, expect, body } =
    Http.request
        { method = "POST"
        , headers = []
        , body = body
        , url = url
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }


put : { url : String, expect : Http.Expect msg, token : String, body : Body } -> Cmd msg
put { url, expect, token, body } =
    Http.request
        { method = "PUT"
        , headers = [ Http.header "Authorization" ("Bearer " ++ token) ]
        , body = body
        , url = url
        , expect = expect
        , timeout = Nothing
        , tracker = Nothing
        }
