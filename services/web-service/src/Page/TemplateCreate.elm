module Page.TemplateCreate exposing (Model, Msg, init, subscriptions, update, view)

import Components
import Html.Styled exposing (Html, div, form, text)
import Html.Styled.Attributes as Attrs
import Html.Styled.Events exposing (onSubmit)
import Http
import Input
import List.Extra
import RemoteData exposing (WebData)
import Session
import Tailwind.Theme as Theme
import Tailwind.Utilities as Tw
import Template.Element as Element


type alias Model =
    { name : String
    , elements : List Element.Form
    , loading : Bool
    , error : Maybe Http.Error
    }


type Msg
    = Submit
    | NameUpdate String
    | ValueFromUpdate Int String
    | XUpdate Int Int
    | YUpdate Int Int
    | WidthUpdate Int Int
    | HeightUpdate Int Int
    | AddElement


init : Session.Token -> ( Model, Cmd Msg )
init token =
    ( { name = ""
      , elements = []
      , loading = False
      , error = Nothing
      }
    , Cmd.none
    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Submit ->
            ( model, Cmd.none )

        NameUpdate value ->
            ( { model | name = value }, Cmd.none )

        AddElement ->
            ( { model | elements = model.elements ++ [ Element.initForm ] }, Cmd.none )

        ValueFromUpdate index value ->
            ( { model
                | elements =
                    List.Extra.updateAt index
                        (\element -> { element | valueFrom = value })
                        model.elements
              }
            , Cmd.none
            )

        XUpdate index value ->
            ( { model
                | elements =
                    List.Extra.updateAt index
                        (\element -> { element | x = value })
                        model.elements
              }
            , Cmd.none
            )

        YUpdate index value ->
            ( { model
                | elements =
                    List.Extra.updateAt index
                        (\element -> { element | y = value })
                        model.elements
              }
            , Cmd.none
            )

        WidthUpdate index value ->
            ( { model
                | elements =
                    List.Extra.updateAt index
                        (\element -> { element | width = value })
                        model.elements
              }
            , Cmd.none
            )

        HeightUpdate index value ->
            ( { model
                | elements =
                    List.Extra.updateAt index
                        (\element -> { element | height = value })
                        model.elements
              }
            , Cmd.none
            )


view : Model -> Html Msg
view model =
    Components.viewContainer "Create Template"
        [ viewForm model ]
        []


viewForm : Model -> Html Msg
viewForm { name, elements, loading } =
    form
        [ Attrs.css
            [ Tw.space_y_6 ]
        , onSubmit Submit
        ]
        [ Input.string
            name
            NameUpdate
            { label = "Name"
            , name = "mail"
            , required = True
            }
        , Input.viewLabel "Elements" "elements" True
        , viewElements elements
        , Components.viewButton [ text "Add" ] "button" False AddElement
        , div [] [ Components.viewSubmitButton loading Submit ]
        ]


viewElements : List Element.Form -> Html Msg
viewElements elements =
        div [ Attrs.css [ Tw.space_y_6 ]] (List.indexedMap viewElement elements)


viewElement : Int -> Element.Form -> Html Msg
viewElement index element =
    div [ Attrs.css [ Tw.space_y_6 ] ]
        [ Input.string
            element.valueFrom
            (ValueFromUpdate index)
            { label = "Property name in JSON"
            , name = "valueFrom"
            , required = True
            }
        , Input.number element.x
            (XUpdate index)
            { label = "X"
            , name = "x"
            , required = True
            }
        , Input.number element.y
            (YUpdate index)
            { label = "Y"
            , name = "y"
            , required = True
            }
        ]


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none
