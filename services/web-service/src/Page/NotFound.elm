module Page.NotFound exposing (view)

import Html.Styled exposing (Html, div, h1, p, text)
import Html.Styled.Attributes as Attrs
import Tailwind.Theme as Theme
import Tailwind.Utilities as Tw


view : Html msg
view =
    div [ Attrs.css [ Tw.min_h_full, Tw.w_full, Tw.mt_64, Tw.mx_auto, Tw.text_center, Tw.text_color Theme.gray_900 ] ]
        [ h1 [ Attrs.css [ Tw.font_extrabold, Tw.text_8xl ] ] [ text "404" ]
        , p [ Attrs.css [ Tw.text_2xl ] ] [ text "Not found" ]
        ]
