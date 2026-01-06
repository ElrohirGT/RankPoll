module Pages.CreatePoll exposing (..)

import Api
import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)
import Http
import Router


type alias Model =
    { title : String
    , options : List String
    , durationInMinutes : Int
    , navigator : Router.Navigator Msg
    , newOption : String
    , error : Maybe String
    }


init : Router.Navigator Msg -> ( Model, Cmd Msg )
init navigator =
    ( { title = ""
      , options = []
      , durationInMinutes = 0
      , navigator = navigator
      , newOption = ""
      , error = Nothing
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
    case msg of
        UpdateTitle newTitle ->
            ( { model | title = newTitle }, Cmd.none )

        UpdateDuration newValidUntil ->
            ( { model | durationInMinutes = newValidUntil }, Cmd.none )

        UpdateNewOption newOption ->
            ( { model | newOption = newOption }, Cmd.none )

        AddNewOption ->
            ( { model
                | options = model.newOption :: model.options
                , newOption = ""
              }
            , Cmd.none
            )

        DeleteOption opt ->
            ( { model
                | options = List.filter (\a -> a /= opt) model.options
              }
            , Cmd.none
            )

        CreatePoll ->
            ( { model | error = Nothing }, Api.createPoll PollCreated model )

        PollCreated res ->
            case res of
                Err _ ->
                    ( { model | error = Just "We're sorry, please try again later." }, Cmd.none )

                Ok resp ->
                    ( model, model.navigator (Router.ViewPoll resp.pollId) )


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


view : Model -> Browser.Document Msg
view model =
    { title = "CreatePoll"
    , body =
        [ input
            [ type_ "text"
            , value model.title
            , onInput UpdateTitle
            ]
            []
            |> displayWithLabel "Title:"
        , input
            [ type_ "number"
            , value (String.fromInt model.durationInMinutes)
            , onInput (\s -> UpdateDuration (Maybe.withDefault 0 (String.toInt s)))
            ]
            []
            |> displayWithLabel "Minutes to vote:"
        , h2 [] [ text "Options:" ]
        , div []
            [ input [ placeholder "New Option:", value model.newOption, onInput UpdateNewOption ] []
            , button [ onClick AddNewOption ] [ text "+" ]
            ]
        , div [] (model.options |> List.map (displayOption DeleteOption))
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
