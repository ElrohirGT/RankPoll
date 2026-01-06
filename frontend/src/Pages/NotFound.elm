module Pages.NotFound exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Router
import Types


type alias Model =
    { user : Types.User
    , navigator : Router.Navigator Msg
    }


init : Router.Navigator Msg -> ( Model, Cmd Msg )
init navigator =
    ( { user =
            { username = ""
            , password = ""
            }
      , navigator = navigator
      }
    , Cmd.none
    )


type Msg
    = ClickedLogin


update : Msg -> Model -> ( Model, Cmd msg )
update msg model =
    ( model, Cmd.none )


view : Model -> Browser.Document Msg
view _ =
    { title = "NotFound!"
    , body =
        []
    }
