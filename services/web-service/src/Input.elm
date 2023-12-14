module Input exposing (Input, InputType(..), view, viewLabeled)

import Css
import Html.Styled exposing (Html, div, input)
import Html.Styled.Attributes as Attrs
import Html.Styled.Events exposing (onInput)
import Tailwind.Breakpoints as Breakpoints
import Tailwind.Theme as Theme
import Tailwind.Utilities as Tw


type InputType
    = Email
    | Text
    | Password


type alias Input a msg =
    { value : a
    , label : String
    , name : String
    , msg : a -> msg
    , required : Bool
    , type_ : InputType
    }



inputTypeToString : InputType -> String
inputTypeToString i =
    case i of
        Email ->
            "email"

        Text ->
            "text"

        Password ->
            "password"


viewLabeled : Input a msg -> Html msg
viewLabeled field =
    div
        []
        [ viewLabel field.label field.name
        , div
            [ Attrs.css
                [ Tw.mt_2 ]
            ]
            [ viewInput field ]
        ]


viewLabel : String -> String -> Html msg
viewLabel value name =
    label
        [ Attrs.css
            [ Tw.block
            , Tw.text_sm
            , Tw.font_medium
            , Tw.text_color Theme.gray_900
            ]
        , Attrs.for name
        ]
        [ text value ]


view : Input a msg -> Html msg
view { value, name, msg, required, type_ } =
    input
        [ Attrs.css
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
        , Attrs.id name
        , Attrs.type_ (inputTypeToString type_)
        , Attrs.required required
        , Attrs.value value
        , onInput msg
        ]
        []
