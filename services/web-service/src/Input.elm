module Input exposing
    ( Input
    , email
    , number
    , password
    , select
    , string
    , view
    , viewLabel
    , viewLabeled
    )

import Css
import Html.Styled exposing (Html, div, label, option, text)
import Html.Styled.Attributes as Attrs
import Html.Styled.Events exposing (onInput)
import Tailwind.Breakpoints as Breakpoints
import Tailwind.Theme as Theme
import Tailwind.Utilities as Tw


type alias Input msg =
    { label : String
    , name : String
    , required : Bool
    , input : InputType msg
    }


type InputType msg
    = Email String (String -> msg)
    | Text String (String -> msg)
    | Password String (String -> msg)
    | Integer Int (Int -> msg)
    | Select String (List ( String, String )) (String -> msg)


type alias InputConfig =
    { label : String
    , name : String
    , required : Bool
    }


email : String -> (String -> msg) -> InputConfig -> Html msg
email value msg config =
    viewLabeled
        { label = config.label
        , name = config.name
        , required = config.required
        , input = Email value msg
        }


string : String -> (String -> msg) -> InputConfig -> Html msg
string value msg config =
    viewLabeled
        { label = config.label
        , name = config.name
        , required = config.required
        , input = Text value msg
        }


password : String -> (String -> msg) -> InputConfig -> Html msg
password value msg config =
    viewLabeled
        { label = config.label
        , name = config.name
        , required = config.required
        , input = Password value msg
        }


number : Int -> (Int -> msg) -> InputConfig -> Html msg
number value msg config =
    viewLabeled
        { label = config.label
        , name = config.name
        , required = config.required
        , input = Integer value msg
        }


select : String -> List ( String, String ) -> (String -> msg) -> InputConfig -> Html msg
select value options msg config =
    viewLabeled
        { label = config.label
        , name = config.name
        , required = config.required
        , input = Select value options msg
        }


viewLabeled : Input msg -> Html msg
viewLabeled input =
    div
        []
        [ viewLabel input.label input.name input.required
        , div
            [ Attrs.css
                [ Tw.mt_2 ]
            ]
            [ view input ]
        ]


viewLabel : String -> String -> Bool -> Html msg
viewLabel value name required =
    let
        labelValue =
            if required then
                value ++ "*"

            else
                value
    in
    label
        [ Attrs.css
            [ Tw.block
            , Tw.text_sm
            , Tw.font_medium
            , Tw.text_color Theme.gray_900
            ]
        , Attrs.for name
        ]
        [ text labelValue ]


view : Input msg -> Html msg
view { name, required, input } =
    case input of
        Email value msg ->
            Html.Styled.input
                [ Attrs.css inputStyle
                , Attrs.id name
                , Attrs.type_ "email"
                , Attrs.required required
                , Attrs.value value
                , onInput msg
                ]
                []

        Text value msg ->
            Html.Styled.input
                [ Attrs.css inputStyle
                , Attrs.id name
                , Attrs.type_ "text"
                , Attrs.required required
                , Attrs.value value
                , onInput msg
                ]
                []

        Password value msg ->
            Html.Styled.input
                [ Attrs.css inputStyle
                , Attrs.id name
                , Attrs.type_ "password"
                , Attrs.required required
                , Attrs.value value
                , onInput msg
                ]
                []

        Integer value msg ->
            Html.Styled.input
                [ Attrs.css inputStyle
                , Attrs.id name
                , Attrs.type_ "number"
                , Attrs.required required
                , Attrs.value (String.fromInt value)
                , onInput
                    (\newValue ->
                        case String.toInt newValue of
                            Just newInt ->
                                msg newInt

                            Nothing ->
                                msg value
                    )
                ]
                []

        Select value options msg ->
            Html.Styled.select
                [ Attrs.css inputStyle
                , Attrs.id name
                , Attrs.required required
                , Attrs.value value
                , onInput msg
                ]
                (List.map (\( label, optionValue ) -> option [ Attrs.value optionValue ] [ text label ]) options)


inputStyle : List Css.Style
inputStyle =
    [ Tw.block
    , Tw.w_full
    , Tw.rounded_md
    , Tw.border_0
    , Tw.py_1_dot_5
    , Tw.px_3
    , Tw.text_color Theme.gray_900
    , Tw.shadow_sm
    , Tw.ring_1
    , Tw.ring_inset
    , Tw.ring_color Theme.gray_300
    , Tw.placeholder_color Theme.gray_400
    , Css.focus
        [ Tw.ring_2
        , Tw.ring_inset
        , Tw.ring_color Theme.indigo_600
        ]
    , Breakpoints.sm
        [ Tw.text_sm
        , Tw.leading_6
        ]
    ]
