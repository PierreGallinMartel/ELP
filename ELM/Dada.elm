module TestHttp exposing (main)

import Browser
import Html exposing (..)
import Html.Events exposing (onClick)
import Http
import Array
import Random

type alias Model =
    { words : List String
     --errorMessage : Maybe String
     , currentOne : String
     , errorMessage : String
     , currentDef : String
     , random : Int
    }

view : Model -> Html Msg
view model =
    div []
        [ 
            button [ onClick SendHttpRequest ][ text "Get data from server" ],
            button [onClick GenerateRandomNumber] [text "generate rand"],
            div [] [text model.errorMessage]
            , viewWord model
        ]

viewRand : Model -> Html Msg
viewRand model=
    makeRandOkay model.random

makeRandOkay : Int -> Html Msg
makeRandOkay random =
    let
        r = String.fromInt random
    in
    div [] [text r]


viewDef : Model -> Html Msg
viewDef model = 
    makeDefOkay model.currentDef


makeDefOkay : String -> Html Msg
makeDefOkay definition =
    div [] [text definition]


viewWord : Model -> Html Msg
viewWord model = 
    makeWordOkay model.currentOne


makeWordOkay : String -> Html Msg
makeWordOkay definition =
    div [] [text definition]

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
    = SendHttpRequest
    | DataReceived (Result Http.Error String)
    | DefReceived (Result Http.Error String)
    | GenerateRandomNumber
    | NewRandomNumber Int


url : String
url =
    "http://localhost:8000/wordList.txt"


getWords : Cmd Msg
getWords =
    Http.get
        { url = url
        , expect = Http.expectString DataReceived
        }

getDef: Cmd Msg
getDef = 
    Http.get
        { url = "http://worldtimeapi.org/api/timezone/Europe/Paris"
        , expect = Http.expectString DefReceived
        }

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GenerateRandomNumber ->
            ( model, Random.generate NewRandomNumber (Random.int 0 100) )

        NewRandomNumber number ->
            let
                arr = Array.fromList model.words
                dudu = Array.get number arr
                a = Maybe.withDefault "........." dudu
            in
            ( {model | random = number, currentOne = a}, Cmd.none )
        
        SendHttpRequest ->
            ( model, getDef )

        DataReceived (Ok wordsStr) ->
            let
                words = String.split " " wordsStr
               -- wL = List.map (\word -> text word) words
            in
            ( { model | words = words}, Cmd.none )

        DataReceived (Err httpError) ->
            ( { model
                | errorMessage = "Problem"
              }
            , Cmd.none
            )
        
        DefReceived (Ok res) ->
            ({model | currentDef = res, errorMessage = "received"}, Cmd.none)
        
        DefReceived (Err httpError) ->
            ( { model
                | errorMessage = "Problem with def"
              }
            , Cmd.none
            )




init : () -> ( Model, Cmd Msg )
init _ =
    ( { words = []
      , currentOne = "",
      errorMessage = "None",
      currentDef = ""
      , random = 0
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
