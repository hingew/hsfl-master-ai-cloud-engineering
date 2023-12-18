port module Session exposing
    ( Session
    , Token
    , authToken
    , authenticated
    , decoder
    , gotToken
    , init
    , navKey
    , setToken
    , tokenValue
    )

import Browser.Navigation as Nav
import Json.Decode as Decode exposing (Decoder)


type Token
    = Token String


type Session
    = LoggedIn Nav.Key Token
    | Guest Nav.Key


init : Nav.Key -> Maybe String -> Session
init key maybeToken =
    case maybeToken of
        Just token ->
            LoggedIn key (Token token)

        Nothing ->
            Guest key


authenticated : Session -> Bool
authenticated session =
    case session of
        LoggedIn _ _ ->
            True

        Guest _ ->
            False


navKey : Session -> Nav.Key
navKey session =
    case session of
        Guest key ->
            key

        LoggedIn key _ ->
            key


authToken : Session -> Maybe Token
authToken session =
    case session of
        Guest _ ->
            Nothing

        LoggedIn _ token ->
            Just token


tokenValue : Token -> String
tokenValue (Token value) =
    value


decoder : Decoder Token
decoder =
    Decode.map Token (Decode.field "access_token" Decode.string)



-- Store the token in the local storage


port storeToken : String -> Cmd msg


setToken : Token -> Cmd msg
setToken (Token value) =
    storeToken value


port getToken : () -> Cmd msg


port gotToken : (String -> msg) -> Sub msg
