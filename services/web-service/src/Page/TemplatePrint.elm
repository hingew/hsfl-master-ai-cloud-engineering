module Page.TemplatePrint exposing (Model, Msg, init, subscriptions, update, view)

import Browser.Navigation as Nav
import Bytes exposing (Bytes)
import Components
import Dict exposing (Dict)
import File.Download as Download
import Html.Styled exposing (Html, div, form, pre, text)
import Html.Styled.Attributes as Attrs
import Html.Styled.Events exposing (onSubmit)
import Http.Extra
import Input
import RemoteData exposing (WebData)
import Route
import Session
import Tailwind.Utilities as Tw
import Template exposing (PrintForm, Template, TemplateId)


type alias Model =
    { form : WebData Template.PrintForm
    , result : WebData Bytes
    }


type Msg
    = Submit
    | GotSubmitResult (WebData Bytes)
    | ValueUpdate String String
    | GotTemplate (WebData Template)


init : TemplateId -> Session.Token -> ( Model, Cmd Msg )
init templateId token =
    ( { form = RemoteData.Loading
      , result = RemoteData.NotAsked
      }
    , Template.fetch templateId token GotTemplate
    )


update : Session.Token -> Nav.Key -> Msg -> Model -> ( Model, Cmd Msg )
update token navKey msg model =
    case msg of
        Submit ->
            case model.form of
                RemoteData.Success form ->
                    ( { model | result = RemoteData.Loading }
                    , Template.print token form GotSubmitResult
                    )

                _ ->
                    ( model, Cmd.none )

        GotTemplate response ->
            ( { model | form = RemoteData.map Template.toForm response }, Cmd.none )

        GotSubmitResult response ->
            case response of
                RemoteData.Success content ->
                    ( model, saveFile content )

                _ ->
                    ( { model | result = response }, Cmd.none )

        ValueUpdate key value ->
            ( { model | form = RemoteData.map (Template.setFormValue key value) model.form }, Cmd.none )


view : Model -> Html Msg
view model =
    Components.viewContainer "Create Template"
        [ viewResult model.result
        , Components.viewRemoteData (viewForm (RemoteData.isLoading model.result)) model.form
        ]
        []


viewResult : WebData Bytes -> Html Msg
viewResult response =
    case response of
        RemoteData.Failure err ->
            Components.viewHttpError (Just err)

        _ ->
            text ""


viewForm : Bool -> PrintForm -> Html Msg
viewForm resultLoading printForm =
    form
        [ Attrs.css
            [ Tw.space_y_6 ]
        , onSubmit Submit
        ]
        (viewFormFields printForm.values
            ++ [ div
                    [ Attrs.css
                        [ Tw.flex
                        , Tw.justify_between
                        ]
                    ]
                    [ Components.viewSubmitButton resultLoading
                    , Components.viewCancelButton Route.TemplateList
                    ]
               ]
        )


viewFormFields : Dict String String -> List (Html Msg)
viewFormFields values =
    values
        |> Dict.toList
        |> List.map viewFormValue


viewFormValue : ( String, String ) -> Html Msg
viewFormValue ( key, value ) =
    Input.string
        value
        (ValueUpdate key)
        { label = key
        , name = key
        , required = True
        }


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none



-- Utilities


saveFile : Bytes -> Cmd Msg
saveFile content =
    Download.bytes "print.pdf" "application/pdf" content
