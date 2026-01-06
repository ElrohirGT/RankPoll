module Pages.CreatePoll exposing (..)

import Api
import Browser
import Dict
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Http
import Router
import Types


type alias Model =
    { room : Types.Room
    , navigator : Router.Navigator Msg
    , newOption : String
    }


init : Router.Navigator Msg -> ( Model, Cmd Msg )
init navigator =
    ( { room =
            { title = ""
            , options = []
            , votes = Dict.empty
            , durationInMinutes = 0
            , summary = Nothing
            }
      , navigator = navigator
      , newOption = ""
      }
    , Cmd.none
    )


type Msg
    = UpdateTitle String
    | UpdateDuration Int
    | UpdateNewOption String
    | AddNewOption
    | CreatePoll
    | DeleteOption String
    | PollCreated (Result Http.Error Api.CreatePollResponse)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    let
        room =
            model.room
    in
    case msg of
        UpdateTitle newTitle ->
            ( { model | room = { room | title = newTitle } }, Cmd.none )

        UpdateDuration newValidUntil ->
            ( { model | room = { room | durationInMinutes = newValidUntil } }, Cmd.none )

        UpdateNewOption newOption ->
            ( { model | newOption = newOption }, Cmd.none )

        AddNewOption ->
            ( { model
                | room = { room | options = model.newOption :: room.options }
                , newOption = ""
              }
            , Cmd.none
            )

        DeleteOption opt ->
            ( { model
                | room = { room | options = List.filter (\a -> a /= opt) room.options }
              }
            , Cmd.none
            )

        CreatePoll ->
            ( model, Api.createPoll PollCreated model.room )

        PollCreated res ->
            ( model, Cmd.none )


view : Model -> Browser.Document Msg
view model =
    { title = "CreatePoll"
    , body =
        [ input
            [ type_ "text"
            , value model.room.title
            , onInput UpdateTitle
            ]
            []
            |> displayWithLabel "Title:"
        , input
            [ type_ "number"
            , value (String.fromInt model.room.durationInMinutes)
            , onInput (\s -> UpdateDuration (Maybe.withDefault 0 (String.toInt s)))
            ]
            []
            |> displayWithLabel "Minutes to vote:"
        , h2 [] [ text "Options:" ]
        , div []
            [ input [ placeholder "New Option:", value model.newOption, onInput UpdateNewOption ] []
            , button [ onClick AddNewOption ] [ text "+" ]
            ]
        , div [] (model.room.options |> List.map (displayOption DeleteOption))
        , button [ onClick CreatePoll ] [ text "Create" ]
        ]
    }


displayWithLabel : String -> Html msg -> Html msg
displayWithLabel inputLabel element =
    div []
        [ label [] [ text inputLabel ]
        , element
        ]


displayOption : (String -> msg) -> String -> Html msg
displayOption onDelete opt =
    div []
        [ text opt
        , button [ onClick (onDelete opt) ] [ text "Delete" ]
        ]
