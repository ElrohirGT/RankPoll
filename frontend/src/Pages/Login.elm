module Pages.Login exposing (..)

import Api
import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Http
import Router
import Types


type alias Model =
    { user : Types.User
    , navigator : Router.Navigator Msg
    , error : Maybe String
    }


init : Router.Navigator Msg -> ( Model, Cmd Msg )
init navigator =
    ( { user =
            { username = ""
            , password = ""
            }
      , navigator = navigator
      , error = Nothing
      }
    , Cmd.none
    )


type Msg
    = ClickedLogin
    | LoginCompleted (Result Http.Error Api.LoginOrCreateUserResponse)
    | UserChanged String
    | PasswordChanged String


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UserChanged username ->
            let
                user =
                    model.user

                newUser =
                    { user | username = username }
            in
            ( { model | user = newUser }, Cmd.none )

        PasswordChanged password ->
            let
                user =
                    model.user

                newUser =
                    { user | password = password }
            in
            ( { model | user = newUser }, Cmd.none )

        ClickedLogin ->
            ( model, Api.loginOrCreateUser LoginCompleted model.user )

        LoginCompleted result ->
            case result of
                Err e ->
                    let
                        a =
                            Debug.log "ERROR:" e
                    in
                    ( { model | error = Just "Failed to login!" }, Cmd.none )

                Ok _ ->
                    -- Later we should probably do something with the response
                    -- Like for example saving a token
                    ( model, model.navigator Router.CreatePoll )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


view : Model -> Browser.Document Msg
view model =
    { title = "Login"
    , body =
        [ input
            [ placeholder "Username"
            , onInput UserChanged
            , value model.user.username
            ]
            []
        , input
            [ placeholder "Password"
            , onInput PasswordChanged
            , value model.user.password
            , type_ "password"
            ]
            []
        , button [ onClick ClickedLogin ] [ text "Login" ]
        ]
    }
