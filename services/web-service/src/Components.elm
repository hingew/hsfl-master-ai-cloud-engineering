module Components exposing
    ( viewButton
    , viewCancelButton
    , viewContainer
    , viewHttpError
    , viewLinkButton
    , viewRemoteData
    , viewSubmitButton
    , viewTable
    )

import Css exposing (disabled)
import Html.Styled
    exposing
        ( Html
        , a
        , button
        , div
        , h1
        , h2
        , header
        , input
        , main_
        , p
        , pre
        , span
        , table
        , tbody
        , td
        , text
        , th
        , thead
        , tr
        )
import Html.Styled.Attributes as Attrs
import Html.Styled.Events exposing (onClick)
import Http
import Http.Extra
import RemoteData exposing (WebData)
import Route
import Svg.Styled exposing (path, svg)
import Svg.Styled.Attributes as SvgAttrs
import Tailwind.Breakpoints as Breakpoint
import Tailwind.Theme as Theme
import Tailwind.Utilities as Tw


viewRemoteData : (a -> Html msg) -> WebData a -> Html msg
viewRemoteData fn response =
    case response of
        RemoteData.Loading ->
            div [] [ text "Loading..." ]

        RemoteData.NotAsked ->
            div [] [ text "Not asked..." ]

        RemoteData.Failure err ->
            pre [] [ text (Http.Extra.errToString err) ]

        RemoteData.Success data ->
            fn data


viewSubmitButton : Bool -> Html msg
viewSubmitButton loading =
    input
        [ Attrs.css (buttonStyles Theme.indigo_600 Theme.indigo_500)
        , Attrs.type_ "submit"
        , Attrs.disabled loading
        , Attrs.value "Submit"
        ]
        [ if loading then
            viewLoading

          else
            text "Submit"
        ]


viewCancelButton : Route.Route -> Html msg
viewCancelButton route =
    a
        [ Attrs.css (buttonStyles Theme.red_600 Theme.red_500)
        , Attrs.href (Route.path route)
        ]
        [ text "Cancel" ]


viewButton : List (Html msg) -> String -> Bool -> msg -> Html msg
viewButton content type_ disabled msg =
    button
        [ Attrs.css (buttonStyles Theme.indigo_600 Theme.indigo_500)
        , onClick msg
        , Attrs.type_ type_
        , Attrs.disabled disabled
        ]
        content


viewLinkButton : String -> Route.Route -> Html msg
viewLinkButton label route =
    a
        [ Attrs.css (buttonStyles Theme.indigo_600 Theme.indigo_500)
        , Attrs.href (Route.path route)
        ]
        [ text label ]


viewLoading : Html msg
viewLoading =
    div
        [ Attrs.attribute "role" "status"
        ]
        [ svg
            [ Attrs.attribute "aria-hidden" "true"
            , SvgAttrs.css
                [ Tw.w_8
                , Tw.h_8
                , Tw.text_color Theme.gray_200
                , Tw.animate_spin
                , Tw.fill_color Theme.blue_600
                ]
            , SvgAttrs.viewBox "0 0 100 101"
            , SvgAttrs.fill "none"
            ]
            [ path
                [ SvgAttrs.d "M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
                , SvgAttrs.fill "currentColor"
                ]
                []
            , path
                [ SvgAttrs.d "M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
                , SvgAttrs.fill "currentFill"
                ]
                []
            ]
        , span
            [ Attrs.css
                [ Tw.sr_only
                ]
            ]
            [ text "Loading..." ]
        ]


viewContainer : String -> List (Html msg) -> List (Html msg) -> Html msg
viewContainer title children actions =
    div
        [ Attrs.css
            [ Tw.min_h_full
            ]
        ]
        [ viewHeader title actions
        , viewContent children
        ]


viewHeader : String -> List (Html msg) -> Html msg
viewHeader title actions =
    header
        [ Attrs.css
            [ Tw.bg_color Theme.white
            , Tw.shadow
            ]
        ]
        [ div
            [ Attrs.css
                [ Tw.flex
                , Tw.justify_between
                , Tw.flex_row
                , Tw.mx_auto
                , Tw.max_w_xl
                , Tw.px_4
                , Tw.py_6
                , Breakpoint.lg
                    [ Tw.px_8
                    ]
                , Breakpoint.sm
                    [ Tw.px_6
                    ]
                ]
            ]
            [ h1
                [ Attrs.css
                    [ Tw.text_3xl
                    , Tw.font_bold
                    , Tw.tracking_tight
                    , Tw.text_color Theme.gray_900
                    ]
                ]
                [ text title ]
            , div [] actions
            ]
        ]


viewContent : List (Html msg) -> Html msg
viewContent children =
    main_ []
        [ div
            [ Attrs.css
                [ Tw.mx_auto
                , Tw.max_w_xl
                , Tw.py_6
                , Breakpoint.lg
                    [ Tw.px_8
                    ]
                , Breakpoint.sm
                    [ Tw.px_6
                    ]
                ]
            ]
            children
        ]


viewTable : List String -> List (a -> Html msg) -> List a -> Html msg
viewTable headers cols data =
    case data of
        [] ->
            div [ Attrs.css [ Tw.min_h_full, Tw.w_full, Tw.flex, Tw.pt_32, Tw.flex_col, Tw.justify_center, Tw.items_center ] ]
                [ h2
                    [ Attrs.css
                        [ Tw.text_xl
                        , Tw.font_bold
                        , Tw.text_color Theme.gray_900
                        ]
                    ]
                    [ text "Its empty here ..." ]
                , p
                    [ Attrs.css
                        [ Tw.text_color Theme.gray_900
                        ]
                    ]
                    [ text "Please create a new one, to get started!" ]
                ]

        _ ->
            table [ Attrs.css [ Tw.table_auto ] ]
                [ thead []
                    [ tr []
                        (List.map (\title -> th [] [ text title ]) headers)
                    ]
                , tbody
                    []
                    (List.map (viewTableRow cols) data)
                ]


viewTableRow : List (a -> Html msg) -> a -> Html msg
viewTableRow cols data =
    tr [] (List.map (\col -> td [] [ col data ]) cols)


buttonStyles : Theme.Color -> Theme.Color -> List Css.Style
buttonStyles baseColor hoverColor =
    [ Tw.rounded_md
    , Tw.px_3
    , Tw.py_1_dot_5
    , Tw.text_sm
    , Tw.font_semibold
    , Tw.leading_6
    , Tw.text_color Theme.white
    , Tw.bg_color baseColor
    , Tw.shadow_sm
    , Css.hover [ Tw.bg_color hoverColor ]
    , Css.focus [ Tw.outline, Tw.outline_2, Tw.outline_offset_2, Tw.outline_color baseColor ]
    ]


viewHttpError : Maybe Http.Error -> Html msg
viewHttpError maybeError =
    case maybeError of
        Just err ->
            div
                [ Attrs.css
                    [ Tw.rounded
                    , Tw.p_6
                    , Tw.bg_color Theme.red_500
                    , Tw.text_color Theme.white
                    , Tw.overflow_scroll
                    , Tw.mb_6
                    ]
                ]
                [ pre
                    []
                    [ text (Http.Extra.errToString err) ]
                ]

        Nothing ->
            text ""
