module Api exposing (..)

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
