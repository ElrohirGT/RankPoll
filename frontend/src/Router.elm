module Router exposing (..)

import Browser.Navigation as Nav
import Url
import Url.Builder exposing (relative)


type Page
    = Login
    | CreatePoll
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


parseUrl : Url.Url -> Page
parseUrl url =
    case url.path of
        "/" ->
            Login

        "/poll" ->
            CreatePoll

        _ ->
            NotFound
