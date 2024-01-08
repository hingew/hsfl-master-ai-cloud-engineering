module Main exposing (main)

import Browser exposing (Document)
import Browser.Navigation as Navigation
import Css.Global
import Html.Styled exposing (Html, div)
import Page.Login
import Page.NotFound
import Page.TemplateCreate
import Page.TemplateList
import Page.TemplatePrint
import Route
import Session exposing (Session)
import Tailwind.Utilities as Tw
import Url
import Url.Parser


type Page
    = Login Page.Login.Model
    | TemplateList Page.TemplateList.Model
    | TemplateCreate Page.TemplateCreate.Model
    | TemplatePrint Page.TemplatePrint.Model
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
    | TemplateCreateMsg Page.TemplateCreate.Msg
    | TemplatePrintMsg Page.TemplatePrint.Msg


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

                Route.Register ->
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

                Route.TemplateCreate ->
                    let
                        ( m, cmd ) =
                            Page.TemplateCreate.init token
                    in
                    ( { model | page = TemplateCreate m }, Cmd.map TemplateCreateMsg cmd )

                Route.TemplatePrint templateId ->
                    let
                        ( m, cmd ) =
                            Page.TemplatePrint.init templateId token
                    in
                    ( { model | page = TemplatePrint m }, Cmd.map TemplatePrintMsg cmd )

                Route.NotFound ->
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
        ( LinkClicked urlRequest, _ ) ->
            case urlRequest of
                Browser.Internal url ->
                    let
                        key =
                            model.session
                                |> Session.navKey
                    in
                    ( model, Navigation.pushUrl key (Url.toString url) )

                Browser.External href ->
                    ( model, Navigation.load href )

        ( UrlChanged url, _ ) ->
            navigate url model

        ( LoginMsg loginMsg, Login login ) ->
            let
                ( updatedPage, cmd ) =
                    Page.Login.update loginMsg login
            in
            ( { model | page = Login updatedPage }, Cmd.map LoginMsg cmd )

        ( TemplateListMsg templateListMsg, TemplateList templateList ) ->
            case Session.authToken model.session of
                Just token ->
                    let
                        ( updatePage, cmd ) =
                            Page.TemplateList.update token templateListMsg templateList
                    in
                    ( { model | page = TemplateList updatePage }, Cmd.map TemplateListMsg cmd )

                Nothing ->
                    ( model, Route.replaceUrl (Session.navKey model.session) Route.Login )

        ( TemplateCreateMsg templateCreateMsg, TemplateCreate templateCreate ) ->
            case Session.authToken model.session of
                Just token ->
                    let
                        ( updatePage, cmd ) =
                            Page.TemplateCreate.update token (Session.navKey model.session) templateCreateMsg templateCreate
                    in
                    ( { model | page = TemplateCreate updatePage }, Cmd.map TemplateCreateMsg cmd )

                Nothing ->
                    ( model, Route.replaceUrl (Session.navKey model.session) Route.Login )

        ( TemplatePrintMsg templatePrintMsg, TemplatePrint templatePrint ) ->
            case Session.authToken model.session of
                Just token ->
                    let
                        ( updatePage, cmd ) =
                            Page.TemplatePrint.update token (Session.navKey model.session) templatePrintMsg templatePrint
                    in
                    ( { model | page = TemplatePrint updatePage }, Cmd.map TemplatePrintMsg cmd )

                Nothing ->
                    ( model, Route.replaceUrl (Session.navKey model.session) Route.Login )

        _ ->
            ( model, Cmd.none )


view : Model -> Document Msg
view model =
    { title = "PDF Designer"
    , body =
        List.map Html.Styled.toUnstyled
            [ Css.Global.global Tw.globalStyles
            , div
                []
                [ viewPage model.page
                ]
            ]
    }


viewPage : Page -> Html Msg
viewPage page =
    case page of
        Login pageModel ->
            Html.Styled.map LoginMsg (Page.Login.view pageModel)

        TemplateList pageModel ->
            Html.Styled.map TemplateListMsg (Page.TemplateList.view pageModel)

        TemplateCreate pageModel ->
            Html.Styled.map TemplateCreateMsg (Page.TemplateCreate.view pageModel)

        TemplatePrint pageModel ->
            Html.Styled.map TemplatePrintMsg (Page.TemplatePrint.view pageModel)

        NotFound ->
            Page.NotFound.view


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
