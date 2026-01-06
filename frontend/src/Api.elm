module Api exposing (..)

import Constants
import Http
import Json.Decode as D
import Json.Encode as E
import Types


basePath : String
basePath =
    "http://127.0.0.1:8080"


type alias LoginOrCreateUserRequest =
    Types.User


loginOrCreateUserRequestEncoder : LoginOrCreateUserRequest -> E.Value
loginOrCreateUserRequestEncoder a =
    E.object
        [ ( "Username", E.string a.username )
        , ( "Password", E.string a.password )
        ]


type alias LoginOrCreateUserResponse =
    { message : String }


loginOrCreateUserResponseDecoder : D.Decoder LoginOrCreateUserResponse
loginOrCreateUserResponseDecoder =
    D.map LoginOrCreateUserResponse (D.field "Msg" D.string)


loginOrCreateUser : (Result Http.Error LoginOrCreateUserResponse -> msg) -> LoginOrCreateUserRequest -> Cmd msg
loginOrCreateUser toMsg req =
    Http.post
        { url = String.concat [ basePath, "/api/user" ]
        , body = loginOrCreateUserRequestEncoder req |> Http.jsonBody
        , expect = Http.expectJson toMsg loginOrCreateUserResponseDecoder
        }


type alias CreatePollRequest a =
    { a
        | title : String
        , options : List String
        , durationInMinutes : Int
    }


createPollRequestEncoder : CreatePollRequest a -> E.Value
createPollRequestEncoder a =
    E.object
        [ ( "Title", E.string a.title )
        , ( "PollOptions", E.list E.string a.options )
        , ( "PollingDuration", E.int (a.durationInMinutes * 60 * Constants.oneSecondInGo) )
        ]


type alias CreatePollResponse =
    { msg : String
    , pollId : String
    }


createPollResponseDecoder : D.Decoder CreatePollResponse
createPollResponseDecoder =
    D.map2 CreatePollResponse
        (D.field "Msg" D.string)
        (D.field "PollId" D.string)


createPoll : (Result Http.Error CreatePollResponse -> msg) -> CreatePollRequest a -> Cmd msg
createPoll toMsg req =
    Http.post
        { url = String.concat [ basePath, "/api/poll" ]
        , body = createPollRequestEncoder req |> Http.jsonBody
        , expect = Http.expectJson toMsg createPollResponseDecoder
        }
