module Time.Extra exposing (toString)

import Time exposing (Posix)



-- 20.4.2020 20:32


toString : Posix -> String
toString posix =
    String.fromInt (Time.toDay Time.utc posix)
        ++ "."
        ++ toMonth posix
        ++ "."
        ++ String.fromInt (Time.toYear Time.utc posix)
        ++ " "
        ++ String.fromInt (Time.toHour Time.utc posix)
        ++ ":"
        ++ String.fromInt (Time.toMinute Time.utc posix)



toMonth : Posix -> String
toMonth posix =
    case Time.toMonth Time.utc posix of
        Time.Jan ->
            "1"

        Time.Feb ->
            "2"

        Time.Mar ->
            "3"

        Time.Apr ->
            "4"

        Time.May ->
            "5"

        Time.Jun ->
            "6"

        Time.Jul ->
            "7"

        Time.Aug ->
            "8"

        Time.Sep ->
            "9"

        Time.Oct ->
            "10"

        Time.Nov ->
            "11"

        Time.Dec ->
            "11"
