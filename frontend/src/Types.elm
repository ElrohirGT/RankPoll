module Types exposing (..)


type alias Flags =
    { user : Maybe User
    }


type alias User =
    { username : String
    , password : String
    }
