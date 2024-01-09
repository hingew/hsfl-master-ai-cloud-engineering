module Page.Login exposing (Model, Msg, init, subscriptions, update, view)

import Auth
import Components
import Css
import Html.Styled exposing (Html, a, div, form, h2, p, text)
import Html.Styled.Attributes as Attrs
import Html.Styled.Events exposing (onSubmit)
import Http
import Input
import Platform.Cmd as Cmd
import RemoteData exposing (WebData)
import Route
import Session exposing (Session)
import Svg.Styled.Events exposing (onClick)
import Tailwind.Theme as Theme
import Tailwind.Utilities as Tw
import Html.Styled exposing (button)



-- https://tailwindui.com/components/application-ui/forms/sign-in-forms


type alias Model =
    { form : Form
    , error : Maybe Http.Error
    , loading : Bool
    , session : Session
    }


type Form
    = LoginForm Auth.Login
    | RegisterForm Auth.Register


type Msg
    = NameUpdate String
    | PasswordUpdate String
    | PasswordConfirmationUpdate String
    | ToggleRegisterLogin
    | Submit
    | GotLoginResult (WebData Session.Token)
    | GotRegisterResult (Result Http.Error ())
    | SessionSaved


init : Session -> ( Model, Cmd Msg )
init session =
    let
        model =
            { form = LoginForm { email = "", password = "" }
            , error = Nothing
            , loading = False
            , session = session
            }
    in
    if Session.authenticated session then
        ( model, Route.replaceUrl (Session.navKey session) Route.TemplateList )

    else
        ( model, Cmd.none )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NameUpdate value ->
            case model.form of
                LoginForm form ->
                    ( { model | form = LoginForm { form | email = value } }, Cmd.none )

                RegisterForm form ->
                    ( { model | form = RegisterForm { form | email = value } }, Cmd.none )

        PasswordUpdate value ->
            case model.form of
                LoginForm form ->
                    ( { model | form = LoginForm { form | password = value } }, Cmd.none )

                RegisterForm form ->
                    ( { model | form = RegisterForm { form | password = value } }, Cmd.none )

        PasswordConfirmationUpdate value ->
            case model.form of
                LoginForm _ ->
                    ( model, Cmd.none )

                RegisterForm form ->
                    ( { model | form = RegisterForm { form | passwordConfirmation = value } }, Cmd.none )

        ToggleRegisterLogin ->
            case model.form of
                LoginForm _ ->
                    ( { model | form = RegisterForm { email = "", password = "", passwordConfirmation = "" } }, Cmd.none )

                RegisterForm _ ->
                    ( { model | form = LoginForm { email = "", password = "" } }, Cmd.none )

        GotLoginResult response ->
            case response of
                RemoteData.Success token ->
                    ( model, Session.setToken token )

                RemoteData.Failure err ->
                    ( { model | error = Just err }, Cmd.none )

                _ ->
                    ( model, Cmd.none )

        GotRegisterResult response ->
            case response of
                Ok _ ->
                    ( { model | form = LoginForm { email = "", password = "" } }, Cmd.none )

                Err err ->
                    ( { model | error = Just err }, Cmd.none )

        Submit ->
            case model.form of
                LoginForm form ->
                    ( model, Auth.login form GotLoginResult )

                RegisterForm form ->
                    ( model, Auth.register form GotRegisterResult )

        SessionSaved ->
            let
                key =
                    Session.navKey model.session
            in
            ( model, Route.replaceUrl key Route.TemplateList )


view : Model -> Html Msg
view model =
    div
        [ Attrs.css
            [ Tw.flex
            , Tw.min_h_full
            , Tw.flex_col
            , Tw.justify_center
            , Tw.px_6
            , Tw.py_12
            ]
        ]
        [ viewHeader
        , viewForm model
        ]


viewForm : Model -> Html Msg
viewForm { form, error, loading } =
    case form of
        LoginForm loginForm ->
            viewLoginForm loginForm error loading

        RegisterForm registerForm ->
            viewRegisterForm registerForm error loading


viewHeader : Html Msg
viewHeader =
    div
        [ Attrs.css
            [ Tw.mx_auto
            , Tw.w_full
            , Tw.max_w_sm
            ]
        ]
        [ h2
            [ Attrs.css [ Tw.mt_10, Tw.text_center, Tw.text_2xl, Tw.leading_9, Tw.tracking_tight, Tw.text_color Theme.gray_900 ] ]
            [ text "Sign in to you account" ]
        ]


viewLoginForm : Auth.Login -> Maybe Http.Error -> Bool -> Html Msg
viewLoginForm loginForm error loading =
    div
        [ Attrs.css
            [ Tw.mt_10
            , Tw.mx_auto
            , Tw.w_full
            , Tw.max_w_sm
            ]
        ]
        [ Components.viewHttpError error
        , form
            [ Attrs.css
                [ Tw.space_y_6 ]
            , onSubmit Submit
            ]
            [ Input.email
                loginForm.email
                NameUpdate
                { label = "Email"
                , name = "mail"
                , required = True
                }
            , Input.password
                loginForm.password
                PasswordUpdate
                { label = "Password"
                , name = "password"
                , required = True
                }
            , div [] [ Components.viewSubmitButton loading]
            ]
        , p [ Attrs.css [ Tw.mt_10, Tw.text_center, Tw.text_sm, Tw.text_color Theme.gray_600 ] ]
            [ text "Not a member? "
            , button
                [ Attrs.css
                    [ Tw.font_semibold
                    , Tw.leading_6
                    , Tw.text_color Theme.blue_600
                    , Css.hover [ Tw.text_color Theme.blue_500, Tw.cursor_pointer ]
                    ]
                , onClick ToggleRegisterLogin
                ]
                [ text "Register now!" ]
            ]
        ]


viewRegisterForm : Auth.Register -> Maybe Http.Error -> Bool -> Html Msg
viewRegisterForm registerForm err loading =
    div
        [ Attrs.css
            [ Tw.mt_10
            , Tw.mx_auto
            , Tw.w_full
            , Tw.max_w_sm
            ]
        ]
        [ Components.viewHttpError err
        , form
            [ Attrs.css
                [ Tw.space_y_6 ]
            , onSubmit Submit
            ]
            [ Input.email
                registerForm.email
                NameUpdate
                { label = "Email"
                , name = "mail"
                , required = True
                }
            , Input.password
                registerForm.password
                PasswordUpdate
                { label = "Password"
                , name = "password"
                , required = True
                }
            , Input.password
                registerForm.passwordConfirmation
                PasswordConfirmationUpdate
                { label = "Password confirmation"
                , name = "password_confirmation"
                , required = True
                }
            , div [] [ Components.viewSubmitButton loading]
            ]
        , p [ Attrs.css [ Tw.mt_10, Tw.text_center, Tw.text_sm, Tw.text_color Theme.gray_600 ] ]
            [ text "Allready a member? "
            , button
                [ Attrs.css
                    [ Tw.font_semibold
                    , Tw.leading_6
                    , Tw.text_color Theme.blue_600
                    , Css.hover [ Tw.text_color Theme.blue_500, Tw.cursor_pointer ]
                    ]
                , onClick ToggleRegisterLogin
                ]
                [ text "Login!" ]
            ]
        ]



subscriptions : Model -> Sub Msg
subscriptions _ =
    Session.gotToken (\_ -> SessionSaved)
