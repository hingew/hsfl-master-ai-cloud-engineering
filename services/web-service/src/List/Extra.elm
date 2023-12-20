module List.Extra exposing (updateAt)


updateAt : Int -> (a -> a) -> List a -> List a
updateAt index fn list =
    case list of
        [] ->
            []

        x :: xs ->
            if index == 0 then
                fn x :: xs

            else
                updateAt (index - 1) fn xs
