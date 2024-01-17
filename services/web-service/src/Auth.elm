module Auth exposing (Login, Register, login, register)

import Http
import Json.Encode as Encode
import RemoteData exposing (WebData)
import Session exposing (Token)


type alias Login =
    { email : String, password : String }


type alias Register =
    { email : String, password : String, passwordConfirmation : String }


encodeLogin : Login -> Encode.Value
encodeLogin { email, password } =
    Encode.object
        [ ( "email", Encode.string email )
        , ( "password", Encode.string password )
        ]


encodeRegister : Register -> Encode.Value
encodeRegister { email, password, passwordConfirmation } =
    Encode.object
        [ ( "email", Encode.string email )
        , ( "password", Encode.string password )
        , ( "password_confirmation", Encode.string passwordConfirmation )
        ]


login : Login -> (WebData Token -> msg) -> Cmd msg
login data msg =
    Http.post
        { url = "/auth/login"
        , expect = Http.expectJson (RemoteData.fromResult >> msg) Session.decoder
        , body = Http.jsonBody (encodeLogin data)
        }


register : Register -> (Result Http.Error () -> msg) -> Cmd msg
register data msg =
    Http.post
        { url = "/auth/register"
        , expect = Http.expectWhatever msg
        , body = Http.jsonBody (encodeRegister data)
        }
