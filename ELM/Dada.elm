module Dada exposing(..)
import Browser
import Html exposing (Html, text, pre)
import Http
import List exposing (..)
import Array
import Debug
-- MAIN


main =
  Browser.element
    { init = init
    , update = update
    , subscriptions = subscriptions
    , view = view
    }



-- MODEL


type Model
  = Failure
  | Loading
  | Success String


init : () -> (Model, Cmd Msg)
init _ =
  ( Loading
  , Http.get
      { url = "http://localhost:8000/wordList.txt"
      , expect = Http.expectString GotText
      }
  )



-- UPDATE


type Msg
  = GotText (Result Http.Error String)
  | GotDef (Result Hrrp.Error String)



update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
  case msg of
    GotText result ->
      case result of
        Ok fullText ->
          (Success fullText, Cmd.none)

        Err _ ->
          (Failure, Cmd.none)

    GotDef result ->
       case result of
        Ok def ->
          (DefinitonOk def, Cmd.none)

        Err _ ->
          (Failure, Cmd.none)



-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none


-- VIEW


view : Model -> Html Msg
view model =
  case model of
    Failure ->
      text "I was unable to load your book."

    Loading ->
      text "Loading..."

    Success fullText ->
      let 
        t = String.split " " fullText
        wL = List.map (\word -> text word) t
        arr = Array.fromList t
        dudu = Array.get 56 arr
        a = Maybe.withDefault "........." dudu
      in
      pre [] [text a]

get n xs  = List.head (List.drop n xs)

