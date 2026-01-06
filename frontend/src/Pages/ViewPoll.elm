module Pages.ViewPoll exposing (..)

import Api
import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Http
import Router
import Task
import Time
import Types


type State
    = Loading
    | HasPoll Types.Room


type alias Model =
    { pollId : String
    , navigator : Router.Navigator Msg
    , error : Maybe String
    , currentTime : Time.Posix
    , state : State
    }


init : String -> Router.Navigator Msg -> ( Model, Cmd Msg )
init pollId navigator =
    ( { pollId = pollId
      , navigator = navigator
      , error = Nothing
      , state = Loading
      , currentTime = Time.millisToPosix 0
      }
    , Cmd.batch
        [ Task.perform TickClock Time.now
        , Api.getPoll GotPoll pollId
        ]
    )


type Msg
    = GotPoll (Result Http.Error Api.GetPollResponse)
    | TickClock Time.Posix


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GotPoll res ->
            case res of
                Err e ->
                    let
                        _ =
                            Debug.log "ERROR:" e
                    in
                    ( { model | error = Just "Failed to get room!" }, Cmd.none )

                Ok response ->
                    ( { model | state = HasPoll response }, Cmd.none )

        TickClock newTime ->
            ( { model | currentTime = newTime }, Cmd.none )


subscriptions : Model -> Sub Msg
subscriptions model =
    case model.state of
        Loading ->
            Sub.none

        HasPoll _ ->
            Time.every 1000 TickClock


view : Model -> Browser.Document Msg
view model =
    { title = "Main Poll View"
    , body =
        case model.state of
            Loading ->
                [ text "Cargando..." ]

            HasPoll room ->
                let
                    remainigMs =
                        Time.posixToMillis room.validUntil - Time.posixToMillis model.currentTime
                in
                h1 [] [ text room.title ] :: displayOptions remainigMs room
    }


displayTimer : Int -> Html msg
displayTimer remainingMs =
    let
        hours =
            remainingMs // (1000 * 60 * 60)

        minutes =
            (remainingMs - hours * 60 * 60 * 1000) // (1000 * 60)

        seconds =
            (remainingMs - (minutes * 60 * 1000 + hours * 60 * 60 * 1000)) // 1000
    in
    div []
        [ text "Remaining:"
        , div []
            [ if hours > 0 then
                text (String.concat [ String.fromInt hours, "hrs " ])

              else
                text ""
            , if minutes > 0 then
                text (String.concat [ String.fromInt minutes, "m " ])

              else
                text ""
            , if seconds > 0 then
                text (String.concat [ String.fromInt seconds, "s" ])

              else
                text ""
            ]
        ]


displayOptions : Int -> Types.Room -> List (Html msg)
displayOptions remainingMs room =
    if remainingMs > 0 then
        [ displayTimer remainingMs
        , h2 [] [ text "Options:" ]
        , div [] (List.map displayOption room.options)
        ]

    else
        [ text "You can no longer vote on this poll!" ]


displayOption : String -> Html msg
displayOption opt =
    div []
        [ text opt
        , button [] [ text "Up" ]
        , button [] [ text "Down" ]
        ]
