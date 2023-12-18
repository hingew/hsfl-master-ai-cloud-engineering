module Main exposing (main)

import Browser exposing (Document)
import Browser.Navigation as Navigation
import Css.Global
import Html.Styled exposing (Html, div, text)
import Page.Login
import Page.TemplateList
import Route
import Session exposing (Session)
import Tailwind.Utilities as Tw
import Url
import Url.Parser


type Page
    = Login Page.Login.Model
    | TemplateList Page.TemplateList.Model
    | NotFound


type alias Model =
    { session : Session
    , page : Page
    }


type Msg
    = UrlChanged Url.Url
    | LinkClicked Browser.UrlRequest
    | LoginMsg Page.Login.Msg
    | TemplateListMsg Page.TemplateList.Msg


type alias Flags =
    { token : Maybe String }


init : Flags -> Url.Url -> Navigation.Key -> ( Model, Cmd Msg )
init flags url key =
    let
        session =
            Session.init key flags.token

        model =
            { session = session, page = NotFound }
    in
    if Session.authenticated session then
        navigate url model

    else
        fromRoute Route.Login model


parseRoute : Url.Url -> Route.Route
parseRoute url =
    case Url.Parser.parse Route.parser url of
        Nothing ->
            Route.NotFound

        Just route ->
            route


fromRoute : Route.Route -> Model -> ( Model, Cmd Msg )
fromRoute route model =
    case Session.authToken model.session of
        Just token ->
            case route of
                Route.Login ->
                    let
                        ( m, cmd ) =
                            Page.Login.init model.session
                    in
                    ( { model | page = Login m }, Cmd.map LoginMsg cmd )

                Route.TemplateList ->
                    let
                        ( m, cmd ) =
                            Page.TemplateList.init token
                    in
                    ( { model | page = TemplateList m }, Cmd.map TemplateListMsg cmd )

                _ ->
                    ( { model | page = NotFound }, Cmd.none )

        Nothing ->
            case model.page of
                Login _ ->
                    ( model, Cmd.none )

                _ ->
                    let
                        ( m, _ ) =
                            Page.Login.init model.session
                    in
                    ( { model | page = Login m }, Route.replaceUrl (Session.navKey model.session) Route.Login )


navigate : Url.Url -> Model -> ( Model, Cmd Msg )
navigate url model =
    fromRoute (parseRoute url) model


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case ( msg, model.page ) of
        ( LoginMsg loginMsg, Login login ) ->
            let
                ( updatedLogin, loginCmd ) =
                    Page.Login.update loginMsg login
            in
            ( { model | page = Login updatedLogin }, Cmd.map LoginMsg loginCmd )

        _ ->
            ( model, Cmd.none )


view : Model -> Document Msg
view model =
    { title = "PDF Designer"
    , body =
        List.map Html.Styled.toUnstyled
            [ div
                []
                [ Css.Global.global Tw.globalStyles
                , viewPage model.page
                ]
            ]
    }


viewPage : Page -> Html Msg
viewPage page =
    case page of
        Login login ->
            Html.Styled.map LoginMsg (Page.Login.view login)

        TemplateList templateList ->
            Html.Styled.map TemplateListMsg (Page.TemplateList.view templateList)

        NotFound ->
            text "Not Found"


subscriptions : Model -> Sub Msg
subscriptions model =
    case model.page of
        Login loginModel ->
            Sub.map LoginMsg (Page.Login.subscriptions loginModel)

        _ ->
            Sub.none


main : Program Flags Model Msg
main =
    Browser.application
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        , onUrlRequest = LinkClicked
        , onUrlChange = UrlChanged
        }
