module Page.TemplateCreate exposing (Model, Msg, init, subscriptions, update, view)

import Browser.Navigation as Nav
import Components
import Html.Styled exposing (Html, div, form, text)
import Html.Styled.Attributes as Attrs
import Html.Styled.Events exposing (onSubmit)
import Http
import Input
import List.Extra
import RemoteData exposing (WebData)
import Route
import Session
import Tailwind.Utilities as Tw
import Template exposing (Template)
import Template.Element as Element


type alias Model =
    { name : String
    , elements : List Element.Form
    , loading : Bool
    , error : Maybe Http.Error
    }


type Msg
    = Submit
    | GotSubmitResult (WebData Template.CreateResponse)
    | NameUpdate String
    | ValueFromUpdate Int String
    | XUpdate Int Int
    | YUpdate Int Int
    | WidthUpdate Int Int
    | HeightUpdate Int Int
    | TypeUpdate Int Element.ElementType
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


update : Session.Token -> Nav.Key -> Msg -> Model -> ( Model, Cmd Msg )
update token navKey msg model =
    case msg of
        Submit ->
            ( { model | loading = True }
            , Template.create token { name = model.name, elements = model.elements } GotSubmitResult
            )

        GotSubmitResult response ->
            case response of
                RemoteData.Success _ ->
                    ( model, Route.replaceUrl navKey Route.TemplateList )

                RemoteData.Failure err ->
                    ( { model | loading = False, error = Just err }, Cmd.none )

                RemoteData.NotAsked ->
                    ( model, Cmd.none )

                RemoteData.Loading ->
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

        TypeUpdate index elementType ->
            ( { model
                | elements =
                    List.Extra.updateAt index
                        (\element -> { element | type_ = elementType })
                        model.elements
              }
            , Cmd.none
            )


view : Model -> Html Msg
view model =
    Components.viewContainer "Create Template"
        [ Components.viewHttpError model.error
        , viewForm model
        ]
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
        , div
            [ Attrs.css
                [ Tw.flex
                , Tw.justify_between
                ]
            ]
            [ Components.viewSubmitButton loading
            , Components.viewCancelButton Route.TemplateList
            ]
        ]


viewElements : List Element.Form -> Html Msg
viewElements elements =
    div [ Attrs.css [ Tw.divide_y ] ] (List.indexedMap viewElement elements)


viewElement : Int -> Element.Form -> Html Msg
viewElement index element =
    div [ Attrs.css [ Tw.space_y_6, Tw.py_6 ] ]
        ([ Input.string
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
         , Input.select
            (Element.typeToString element.type_)
            Element.typeOptions
            (\string ->
                string
                    |> Element.typeFromString
                    |> Maybe.withDefault element.type_
                    |> TypeUpdate index
            )
            { label = "Type"
            , name = "type"
            , required = True
            }
         ]
            ++ viewElementType index element.type_
        )


viewElementType : Int -> Element.ElementType -> List (Html Msg)
viewElementType index type_ =
    case type_ of
        Element.Rect ->
            []

        Element.Text text fontSize ->
            [ Input.string
                text
                (\newValue -> TypeUpdate index (Element.Text newValue fontSize))
                { label = "Text"
                , name = "text"
                , required = True
                }
            , Input.number fontSize
                (\newValue -> TypeUpdate index (Element.Text text newValue))
                { label = "Font size"
                , name = "fontSize"
                , required = True
                }
            ]


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none
