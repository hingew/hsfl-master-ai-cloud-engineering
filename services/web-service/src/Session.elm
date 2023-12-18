port module Session exposing (Session, Token, authToken, authenticated, decoder, init, navKey)

import Browser.Navigation as Nav
import Json.Decode as Decode exposing (Decoder)


type Token
    = Token String


type Session
    = LoggedIn Nav.Key Token
    | Guest Nav.Key


init : Nav.Key -> Maybe String -> ( Session, Cmd msg )
init key maybeToken =
    case maybeToken of
        Just token ->
            ( LoggedIn key (Token token), setToken (Token token) )

        Nothing ->
            ( Guest key, Cmd.none )


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
