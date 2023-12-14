module Main exposing (main)

import Browser
import Css.Global
import Html.Styled exposing (Html, div)
import Page.Login
import Tailwind.Utilities as Tw


type Page
    = Login Page.Login.Model


type alias Model =
    { page : Page }


type Msg
    = HandleLoginMsg Page.Login.Msg


init : ( Model, Cmd Msg )
init =
    let
        ( login, loginCmd ) =
            Page.Login.init
    in
    ( { page = Login login
      }
    , Cmd.map HandleLoginMsg loginCmd
    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case ( msg, model.page ) of
        ( HandleLoginMsg loginMsg, Login login ) ->
            let
                ( updatedLogin, loginCmd ) =
                    Page.Login.update loginMsg login
            in
            ( { model | page = Login updatedLogin }, Cmd.map HandleLoginMsg loginCmd )


view : Model -> Html Msg
view model =
    div
        []
        [ Css.Global.global Tw.globalStyles
        , viewPage model.page
        ]


viewPage : Page -> Html Msg
viewPage page =
    case page of
        Login login ->
            Html.Styled.map HandleLoginMsg (Page.Login.view login)


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none


main : Program () Model Msg
main =
    Browser.element
        { init = \_ -> init
        , view = \model -> Html.Styled.toUnstyled (view model)
        , update = update
        , subscriptions = subscriptions
        }
