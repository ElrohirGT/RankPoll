module Pages.Login exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Types


type alias Model =
    { user : Types.User
    }


init : Types.Flags -> { user : { username : String, password : String } }
init _ =
    { user =
        { username = ""
        , password = ""
        }
    }


type Msg
    = ClickedLogin
    | UserChanged String
    | PasswordChanged String


update : Msg -> Model -> ( Model, Cmd msg )
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
            ( model, Cmd.none )


view : Model -> Browser.Document Msg
view model =
    { title = "Login"
    , body =
        [ input [ placeholder "Username", onInput UserChanged, value model.user.username ] []
        , input [ placeholder "Password", onInput PasswordChanged, value model.user.password ] []
        , button [ onClick ClickedLogin ] [ text "Login" ]
        ]
    }
