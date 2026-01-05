module Main exposing (Model, Msg(..), init, main, subscriptions, update, view)

import Browser
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (..)
import Pages.CreatePoll as CreatePollPage
import Pages.Login as LoginPage
import Types
import Url



-- MAIN


main : Program Types.Flags Model Msg
main =
    Browser.application
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        , onUrlChange = UrlChanged
        , onUrlRequest = LinkClicked
        }



-- MODEL


type Pages
    = LoginView LoginPage.Model
    | CreatePollView CreatePollPage.Model


type alias Model =
    { key : Nav.Key
    , url : Url.Url
    , page : Pages
    }


init : Types.Flags -> Url.Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url key =
    let
        loginModel =
            LoginView (LoginPage.init flags)
    in
    ( Model key url loginModel, Cmd.none )



-- UPDATE


type Msg
    = LinkClicked Browser.UrlRequest
    | UrlChanged Url.Url
    | LoginViewMsg LoginPage.Msg
    | CreatePollViewMsg CreatePollPage.Msg


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    let
        handlePage pageUpdate pageView toMsg innerMsg innerModel =
            let
                res =
                    pageUpdate innerMsg innerModel

                newModel =
                    { model | page = pageView (Tuple.first res) }

                newCmd =
                    Cmd.map toMsg (Tuple.second res)
            in
            ( newModel, newCmd )
    in
    case msg of
        LinkClicked urlRequest ->
            case urlRequest of
                Browser.Internal url ->
                    ( model, Nav.pushUrl model.key (Url.toString url) )

                Browser.External href ->
                    ( model, Nav.load href )

        UrlChanged url ->
            ( { model | url = url }
            , Cmd.none
            )

        LoginViewMsg innerMsg ->
            case model.page of
                LoginView innerModel ->
                    handlePage LoginPage.update LoginView LoginViewMsg innerMsg innerModel

                _ ->
                    ( model, Cmd.none )

        CreatePollViewMsg innerMsg ->
            case model.page of
                CreatePollView innerModel ->
                    handlePage CreatePollPage.update CreatePollView CreatePollViewMsg innerMsg innerModel

                _ ->
                    ( model, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none



-- VIEW


view : Model -> Browser.Document Msg
view model =
    let
        toView innerModel toMsg innerViewFunc =
            let
                innerView =
                    innerViewFunc innerModel

                viewMapper a =
                    Html.map toMsg a
            in
            { title = innerView.title
            , body = List.map viewMapper innerView.body
            }
    in
    case model.page of
        LoginView innerModel ->
            toView innerModel LoginViewMsg LoginPage.view

        CreatePollView innerModel ->
            toView innerModel CreatePollViewMsg CreatePollPage.view
