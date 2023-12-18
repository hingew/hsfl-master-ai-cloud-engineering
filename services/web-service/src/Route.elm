module Route exposing (Route(..), parser, path, title, replaceUrl)

import Browser.Navigation as Nav
import Template exposing (TemplateId)
import Url.Parser as Parser exposing ((</>), Parser)


type Route
    = Login
    | Register
    | TemplateList
    | TemplateDetail TemplateId
    | TemplateCreate
    | TemplateEdit TemplateId
    | NotFound


parser : Parser (Route -> a) a
parser =
    Parser.oneOf
        [ Parser.map TemplateList Parser.top
        , Parser.map TemplateDetail (Parser.s "templates" </> templateIdParser)
        , Parser.map TemplateCreate (Parser.s "templates" </> Parser.s "new")
        , Parser.map TemplateEdit (Parser.s "templates" </> templateIdParser </> Parser.s "edit")
        , Parser.map Login (Parser.s "login")
        , Parser.map Register (Parser.s "register")
        ]


templateIdParser : Parser (TemplateId -> a) a
templateIdParser =
    Parser.map Template.toId Parser.int


title : Route -> String
title route =
    case route of
        Login ->
            "Login"

        Register ->
            "Register"

        TemplateList ->
            "Templates"

        TemplateDetail _ ->
            "Template Detail"

        TemplateCreate ->
            "Create Template"

        TemplateEdit _ ->
            "Edit Template"

        NotFound ->
            "Not found"


path : Route -> String
path route =
    case route of
        Login ->
            "/login"

        Register ->
            "/register"

        TemplateList ->
            "/"

        TemplateDetail id ->
            "/templates/"
                ++ String.fromInt (Template.id id)

        TemplateCreate ->
            "/templates/new"

        TemplateEdit id ->
            "/templates/"
                ++ String.fromInt (Template.id id)
                ++ "/edit"

        NotFound ->
            "/404"


replaceUrl : Nav.Key -> Route -> Cmd msg
replaceUrl key route =
    Nav.replaceUrl key (path route)
