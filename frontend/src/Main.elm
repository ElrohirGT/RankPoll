module Main exposing (Model, Msg(..), init, main, subscriptions, update, view)

import Browser
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (..)
import Pages.CreatePoll as CreatePollPage
import Pages.Login as LoginPage
import Pages.NotFound as NotFoundPage
import Pages.ViewPoll as ViewPollPage
import Router
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
    | ViewPollView ViewPollPage.Model
    | NotFoundView NotFoundPage.Model


type alias Model =
    { key : Nav.Key
    , page : Pages
    }


init : Types.Flags -> Url.Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url key =
    redirectToPage key url



-- UPDATE


type Msg
    = LinkClicked Browser.UrlRequest
    | UrlChanged Url.Url
    | LoginViewMsg LoginPage.Msg
    | CreatePollViewMsg CreatePollPage.Msg
    | NotFoundViewMsg NotFoundPage.Msg
    | ViewPollViewMsg ViewPollPage.Msg


redirectToPage : Nav.Key -> Url.Url -> ( Model, Cmd Msg )
redirectToPage key url =
    let
        newPage =
            Router.fromUrl url
    in
    case newPage of
        Router.Login ->
            let
                ( innerModel, innerCmd ) =
                    key
                        |> Router.createNavigator
                        |> LoginPage.init
            in
            ( Model key (LoginView innerModel)
            , Cmd.map LoginViewMsg innerCmd
            )

        Router.CreatePoll ->
            let
                ( innerModel, innerCmd ) =
                    key
                        |> Router.createNavigator
                        |> CreatePollPage.init
            in
            ( Model key (CreatePollView innerModel)
            , Cmd.map CreatePollViewMsg innerCmd
            )

        Router.NotFound ->
            let
                ( innerModel, innerCmd ) =
                    key
                        |> Router.createNavigator
                        |> NotFoundPage.init
            in
            ( Model key (NotFoundView innerModel)
            , Cmd.map NotFoundViewMsg innerCmd
            )

        Router.ViewPoll id ->
            let
                ( innerModel, innerCmd ) =
                    key
                        |> Router.createNavigator
                        |> ViewPollPage.init id
            in
            ( Model key (ViewPollView innerModel)
            , Cmd.map ViewPollViewMsg innerCmd
            )


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
            redirectToPage model.key url

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

        NotFoundViewMsg innerMsg ->
            case model.page of
                NotFoundView innerModel ->
                    handlePage NotFoundPage.update NotFoundView NotFoundViewMsg innerMsg innerModel

                _ ->
                    ( model, Cmd.none )

        ViewPollViewMsg innerMsg ->
            case model.page of
                ViewPollView innerModel ->
                    handlePage ViewPollPage.update ViewPollView ViewPollViewMsg innerMsg innerModel

                _ ->
                    ( model, Cmd.none )



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
    let
        handleSubs toMsg subs innerModel =
            Sub.map toMsg (subs innerModel)
    in
    case model.page of
        LoginView innerModel ->
            handleSubs LoginViewMsg LoginPage.subscriptions innerModel

        CreatePollView innerModel ->
            handleSubs CreatePollViewMsg CreatePollPage.subscriptions innerModel

        ViewPollView innerModel ->
            handleSubs ViewPollViewMsg ViewPollPage.subscriptions innerModel

        NotFoundView innerModel ->
            handleSubs NotFoundViewMsg NotFoundPage.subscriptions innerModel



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

        ViewPollView innerModel ->
            toView innerModel ViewPollViewMsg ViewPollPage.view

        NotFoundView innerModel ->
            toView innerModel NotFoundViewMsg NotFoundPage.view
