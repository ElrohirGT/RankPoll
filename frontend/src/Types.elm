module Types exposing (..)

import Dict exposing (Dict)
import Json.Decode as D
import Time


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


rankDecoder : D.Decoder Rank
rankDecoder =
    D.map2 Rank
        (D.field "Option" D.string)
        (D.field "Position" D.int)


type alias Vote =
    { username : String
    , ranking : List Rank
    }


voteDecoder : D.Decoder Vote
voteDecoder =
    D.map2 Vote
        (D.field "Username" D.string)
        (D.field "Ranking" (D.list rankDecoder))


type alias Room =
    { title : String
    , options : List String
    , votes : Dict String Vote
    , validUntil : Time.Posix
    , summary : Maybe PollSummary
    }


posixTimeDecoder : D.Decoder Time.Posix
posixTimeDecoder =
    D.map Time.millisToPosix D.int


roomDecoder : D.Decoder Room
roomDecoder =
    D.map5 Room
        (D.field "Title" D.string)
        (D.field "Options" (D.list D.string))
        (D.field "Votes" (D.dict voteDecoder))
        (D.field "ValidUntil" posixTimeDecoder)
        (D.field "Summary" (D.maybe pollSummaryDecoder))


type alias PollSummary =
    { winner : String
    , winnerVoteCount : Int
    , totalVoteCount : Int
    , rounds : List (Dict String Int)
    }


pollSummaryDecoder : D.Decoder PollSummary
pollSummaryDecoder =
    D.map4 PollSummary
        (D.field "Winner" D.string)
        (D.field "WinnerVoteCount" D.int)
        (D.field "TotalVoteCount" D.int)
        (D.field "Rounds" (D.list <| D.dict <| D.int))
