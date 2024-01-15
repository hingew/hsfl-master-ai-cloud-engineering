module List.Extra exposing (updateAt, removeAt)


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


removeAt : Int -> List a -> List a
removeAt index list =
    case list of
        [] ->
            []

        x :: xs ->
            if index == 0 then
                xs

            else
                x :: removeAt (index - 1) xs
