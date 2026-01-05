module Pages.CreatePoll exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
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


update : Msg -> Model -> ( Model, Cmd msg )
update msg model =
    ( model, Cmd.none )


view : Model -> Browser.Document Msg
view _ =
    { title = "CreatePoll"
    , body =
        []
    }
