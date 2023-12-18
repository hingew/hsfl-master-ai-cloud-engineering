module Api exposing
    ( get
    , post
    , put
    )

import Http exposing (Body)


get : { url : String, expect : Http.Expect msg } -> Cmd msg
get { url, expect } =
    Http.request
        { method = "GET"

        --, headers = [ Http.header "Authorization" ("Bearer " ++ token) ]
        , headers = []
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
        , headers []
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
