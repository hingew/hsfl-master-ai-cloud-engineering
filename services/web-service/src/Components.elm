module Components exposing (viewRemoteData, viewLabeldInput)

import Css
import Html.Styled exposing (Html, div, input, label, pre, text)
import Html.Styled.Attributes as Attrs
import Html.Styled.Events exposing (onInput)
import Http.Extra
import RemoteData exposing (WebData)




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






