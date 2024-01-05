module Page.TemplateList exposing (Model, Msg, init, subscriptions, update, view)

import Components
import Html.Styled exposing (Html, div, text)
import Http
import List.Extra
import RemoteData exposing (WebData)
import Route
import Session
import Template exposing (Template, TemplateId)
import Time.Extra


type alias Model =
    { templates : WebData (List Template)
    , error : Maybe Http.Error
    }


type Msg
    = GotTemplates (WebData (List Template))
    | DeleteTemplate TemplateId
    | GotDeleteResponse (WebData ())


init : Session.Token -> ( Model, Cmd Msg )
init token =
    ( { templates = RemoteData.Loading
      , error = Nothing
      }
    , Template.fetchAll token GotTemplates
    )


update : Session.Token -> Msg -> Model -> ( Model, Cmd Msg )
update token msg model =
    case msg of
        GotTemplates response ->
            ( { model | templates = response }, Cmd.none )

        DeleteTemplate id ->
            ( { model | templates = RemoteData.map (\templates -> List.filter (\template -> template.id /= id) templates) model.templates }
            , Template.delete token id GotDeleteResponse
            )

        GotDeleteResponse response ->
            case response of
                RemoteData.Success _ ->
                    ( model, Template.fetchAll token GotTemplates )

                RemoteData.Failure err ->
                    ( { model | error = Just err }, Cmd.none )

                RemoteData.Loading ->
                    ( model, Cmd.none )

                RemoteData.NotAsked ->
                    ( model, Cmd.none )


view : Model -> Html Msg
view model =
    Components.viewContainer "Templates"
        [ Components.viewRemoteData viewTemplates model.templates ]
        [ viewCreate ]


viewCreate : Html Msg
viewCreate =
    Components.viewLinkButton "Create Template" Route.TemplateCreate


viewTemplates : List Template -> Html Msg
viewTemplates templates =
    let
        headers =
            [ "Name", "Created at", "Actions" ]

        cols =
            [ .name >> text
            , .createdAt >> Time.Extra.toString >> text
            , viewActions
            ]
    in
    Components.viewTable headers cols templates


viewActions : Template -> Html Msg
viewActions template =
    div []
        [ Components.viewButton [ text "Delete" ] "button" False (DeleteTemplate template.id) ]


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none
