module Page.TemplateList exposing (Model, Msg, init, subscriptions, update, view)

import Html.Styled exposing (Html, div)
import Session exposing (Session)


type alias Model =
    {}


type Msg
    = NoOp


init : Session -> ( Model, Cmd Msg )
init session =
    ( {}, Cmd.none )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NoOp ->
            ( model, Cmd.none )


view : Model -> Html Msg
view model =
    div [] []


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none
