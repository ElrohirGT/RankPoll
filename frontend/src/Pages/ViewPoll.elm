module Pages.ViewPoll exposing (..)

import Api
import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Http
import Router


type alias Model =
    { pollId : String
    , navigator : Router.Navigator Msg
    , error : Maybe String
    }


init : String -> Router.Navigator Msg -> ( Model, Cmd Msg )
init pollId navigator =
    ( { pollId = pollId
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
    ( model, Cmd.none )


view : Model -> Browser.Document Msg
view model =
    { title = "Main Poll View"
    , body = []
    }
