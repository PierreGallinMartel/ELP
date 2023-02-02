module Projet exposing (main)

import Browser
import Html exposing (..)
import Html.Events exposing (onClick, onInput)
import Html.Attributes exposing (..)
import Http
import Array
import Random
import Json.Decode exposing (Decoder, map, field, int, string, at, list)

type alias Model =
    { words : List String
     , currentOne : String
     , errorMessage : String
     , currentDef : String
     , random : Int
     , guess : String
     , rightOrWrong : String
    }

view : Model -> Html Msg
view model =
    div []
        [ 
            --button [ onClick SendHttpRequest ][ text "Get data from server" ]
              button [onClick GenerateRandomNumber] [text "Get random word"]
            , div[][]
            , input [ placeholder "Guess word", value model.guess, onInput Change] []
            , button [onClick GuessWord] [text "Confirm word"]
            --, div [] [text model.errorMessage]
            --, viewRand model
            , viewRes model
            , div[][]
            , div[][]
            --, viewWord model
            , viewDef model
        ]

viewRand : Model -> Html Msg
viewRand model=
    div[] [text (String.fromInt (model.random))]

viewRes : Model -> Html Msg
viewRes model=
    div[] [text (model.rightOrWrong)]

viewDef : Model -> Html Msg
viewDef model = 
    div[] [text (model.currentDef)]

viewGuess : Model -> Html Msg
viewGuess model = 
    div[] [text (model.guess)]

viewWord : Model-> Html Msg
viewWord model =
    div [] [text (model.currentOne)]

-- viewWord : Model -> Html Msg
-- viewWord model = 
--     makeWordOkay model.currentOne


-- makeWordOkay : String -> Html Msg
-- makeWordOkay w=
--     div [] [text w]

-- viewFullList : Model -> Html Msg
-- viewFullList model =
--     viewNicknames model.words 


-- viewNicknames : List String -> Html Msg
-- viewNicknames words =
--     div []
--         [ h3 [] [ text "Old School Main Characters" ]
--         , ul [] (List.map viewNickname words )
--         ]


-- viewNickname : String -> Html Msg
-- viewNickname nickname =
--     li [] [ text nickname ]


type Msg
    = --SendHttpRequest
     DataReceived (Result Http.Error String)
    | DefReceived (Result Http.Error (List (List (List String))))
    | GenerateRandomNumber
    | NewRandomNumber Int
    | GuessWord
    | Change String


url : String
url =
    "http://localhost:8000/wordList.txt"


-- getWords : Cmd Msg
-- getWords =
--     Http.get
--         { url = url
--         , expect = Http.expectString DataReceived
--         }

-- getDef: Cmd Msg
-- getDef = 
--     Http.get
--         { url = "http://worldtimeapi.org/api/timezone/Europe/Paris"
--         , expect = Http.expectString DefReceived
--         }

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GenerateRandomNumber ->
            ( model, Random.generate NewRandomNumber (Random.int 0 999) )

        NewRandomNumber number ->
            let
                arr = Array.fromList model.words
                dudu = Array.get number arr
                a = Maybe.withDefault "........." dudu
            in
            ( {model | random = number, currentOne = a, rightOrWrong = "Guess the word !"}, Http.get
                { url = "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ a
                , expect = Http.expectJson DefReceived decoder1
                }
            )
        Change newContent ->
            ({ model | guess = newContent }, Cmd.none)
        
        GuessWord ->
            if model.currentOne == "" then ({model |rightOrWrong = "You haven't started the game !", guess = ""}, Cmd.none)
            else if model.guess == "" then ({model |rightOrWrong = "You haven't written anything...", guess=""}, Cmd.none)
            else if model.guess == model.currentOne then ({model |rightOrWrong = ("Yay :) It was " ++ model.currentOne ++ " !"), guess = ""}, Cmd.none)
            else ({model | rightOrWrong = "Try again !", guess = ""}, Cmd.none)

        -- SendHttpRequest ->
        --     ( model, getDef )

        DataReceived (Ok wordsStr) ->
            let
                words = String.split " " wordsStr
            in
            ( { model | words = words}, Cmd.none )

        DataReceived (Err httpError) ->
            ( { model
                | errorMessage = "Problem"
              }
            , Cmd.none
            )
        
        DefReceived (Ok res) ->
            let
                arr = Array.fromList res
                subList = Array.get 0 arr
                c = Maybe.withDefault ([["........."]]) subList
                d = List.map (\def -> String.join " " def) c
            in
            ({model | currentDef = (String.join " " d), errorMessage = "received"}, Cmd.none)
        
        DefReceived (Err httpError) ->
            ( { model
                | errorMessage = "Problem with def"
              }
            , Cmd.none
            )


decoder4 : Decoder String
decoder4 =
    field "definition" string

decoder3 = at ["definitions"] (Json.Decode.list decoder4)

decoder2 = at ["meanings"] (Json.Decode.list decoder3 )

decoder1 = Json.Decode.list decoder2


  

init : () -> ( Model, Cmd Msg )
init _ =
    ( { words = []
      , currentOne = "",
      errorMessage = "None",
      currentDef = ""
      , random = 0
      , guess = ""
      , rightOrWrong = "Start the game !"
      }
    ,   Http.get
        { url = "http://localhost:8000/wordList.txt"
        , expect = Http.expectString DataReceived
        }
    )


main : Program () Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = \_ -> Sub.none
        }
