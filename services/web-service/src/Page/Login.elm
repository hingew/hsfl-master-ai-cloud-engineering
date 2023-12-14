module Page.Login exposing (Model, Msg, init, update, view)

import Components
import Html.Styled exposing (Html, div, form, h2, text)
import Html.Styled.Attributes as Attrs
import Platform.Cmd as Cmd
import Tailwind.Theme as Theme
import Tailwind.Utilities as Tw



-- https://tailwindui.com/components/application-ui/forms/sign-in-forms


type Model
    = Login LoginForm
    | Register RegisterForm


type alias LoginForm =
    { name : String
    , password : String
    }


type alias RegisterForm =
    { name : String
    , password : String
    , passwordConfirmation : String
    }


type Msg
    = NameUpdate String
    | PasswordUpdate String
    | PasswordConfirmationUpdate String
    | ToggleRegisterLogin
    | Submit


init : ( Model, Cmd Msg )
init =
    ( Login { name = "", password = "" }, Cmd.none )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case model of
        Login form ->
            updateLoginForm msg form

        Register form ->
            updateRegisterForm msg form


updateLoginForm : Msg -> LoginForm -> ( Model, Cmd Msg )
updateLoginForm msg form =
    case msg of
        NameUpdate value ->
            ( Login { form | name = value }, Cmd.none )

        PasswordUpdate value ->
            ( Login { form | password = value }, Cmd.none )

        PasswordConfirmationUpdate _ ->
            ( Login form, Cmd.none )

        ToggleRegisterLogin ->
            ( Register { name = "", password = "", passwordConfirmation = "" }, Cmd.none )

        Submit ->
            -- TODO: Handle submit
            ( Login form, Cmd.none )


updateRegisterForm : Msg -> RegisterForm -> ( Model, Cmd Msg )
updateRegisterForm msg form =
    case msg of
        NameUpdate value ->
            ( Register { form | name = value }, Cmd.none )

        PasswordUpdate value ->
            ( Register { form | password = value }, Cmd.none )

        PasswordConfirmationUpdate _ ->
            ( Register form, Cmd.none )

        ToggleRegisterLogin ->
            ( Login { name = "", password = "" }, Cmd.none )

        Submit ->
            -- TODO: Handle submit
            ( Register form, Cmd.none )


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
viewForm model =
    div
        [ Attrs.css
            [ Tw.mt_10
            , Tw.mx_auto
            , Tw.w_full
            , Tw.max_w_sm
            ]
        ]
        [ case model of
            Login form ->
                viewLoginForm form

            Register form ->
                viewRegisterform form
        ]


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


viewLoginForm : LoginForm -> Html Msg
viewLoginForm loginForm =
    form
        [ Attrs.css
            [ Tw.space_y_6 ]
        ]
        [ Components.viewLabeldInput
            { value = loginForm.name
            , label = "Username"
            , name = "username"
            , msg = NameUpdate
            , required = True
            , type_ = Email
            }
        ]


viewRegisterform : RegisterForm -> Html Msg
viewRegisterform form =
    div [] []
