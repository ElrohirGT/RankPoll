module Router exposing (..)

import Browser.Navigation as Nav
import Url
import Url.Builder exposing (relative)
import Url.Parser as P exposing (..)


type Page
    = Login
    | CreatePoll
    | ViewPoll String
    | NotFound


type alias Navigator msg =
    Page -> Cmd msg


createNavigator : Nav.Key -> Page -> Cmd msg
createNavigator key =
    \page -> page |> toString |> Nav.pushUrl key


toString : Page -> String
toString page =
    case page of
        Login ->
            relative [ "" ] []

        CreatePoll ->
            relative [ "poll" ] []

        NotFound ->
            relative [ "404" ] []

        ViewPoll id ->
            relative [ "poll", id ] []


urlParser : Parser (Page -> c) c
urlParser =
    P.oneOf
        [ P.map Login P.top
        , P.map CreatePoll (P.s "poll")
        , P.map ViewPoll (P.s "poll" </> P.string)
        ]


fromUrl : Url.Url -> Page
fromUrl url =
    Maybe.withDefault NotFound (P.parse urlParser url)
