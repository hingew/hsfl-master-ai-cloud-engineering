module Page.TemplateList exposing (Model, Msg, init, subscriptions, update, view)

import Components
import Html.Styled exposing (Html, div, text)
import RemoteData exposing (WebData)
import Session
import Template exposing (Template)
import Time.Extra
import Route


type alias Model =
    { templates : WebData (List Template) }


type Msg
    = GotTemplates (WebData (List Template))


init : Session.Token -> ( Model, Cmd Msg )
init token =
    ( { templates = RemoteData.Loading }, Template.fetchAll token GotTemplates )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GotTemplates response ->
            ( { model | templates = response }, Cmd.none )


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


viewActions : Template -> Html msg
viewActions template =
    div [] [ text "Actions" ]


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none
