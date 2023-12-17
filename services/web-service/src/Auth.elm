module Auth exposing (Login, Register, Token, login, register)

import Http
import Json.Decode as Decode exposing (Decoder)
import Json.Encode as Encode
import RemoteData exposing (WebData)


type Token
    = Token String


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
encodeRegister { email, password } =
    Encode.object
        [ ( "email", Encode.string email )
        , ( "password", Encode.string password )
        ]


tokenDecoder : Decoder Token
tokenDecoder =
    Decode.map Token (Decode.field "access_token" Decode.string)


login : Login -> (WebData Token -> msg) -> Cmd msg
login data msg =
    Http.post
        { url = "/auth/login"
        , expect = Http.expectJson (RemoteData.fromResult >> msg) tokenDecoder
        , body = Http.jsonBody (encodeLogin data)
        }


register : Register -> (Result Http.Error () -> msg) -> Cmd msg
register data msg =
    Http.post
        { url = "/auth/register"
        , expect = Http.expectWhatever msg
        , body = Http.jsonBody (encodeRegister data)
        }
