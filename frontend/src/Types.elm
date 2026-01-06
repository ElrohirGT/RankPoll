module Types exposing (..)

import Dict exposing (Dict)


type alias Flags =
    { user : Maybe User
    }


type alias User =
    { username : String
    , password : String
    }


type alias Rank =
    { option : String
    , position : Int
    }


type alias Vote =
    { username : String
    , ranking : List Rank
    }


type alias Room =
    { title : String
    , options : List String
    , votes : Dict String Vote
    , durationInMinutes : Int
    , summary : Maybe PollSummary
    }


type alias PollSummary =
    { winner : String
    , winnerVoteCount : Int
    , totalVoteCount : Int
    , rounds : List (Dict String Int)
    }
