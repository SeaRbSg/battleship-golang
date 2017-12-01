package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dan-v/seattlerb-battleship/battleship"
	"github.com/dan-v/seattlerb-battleship/common"
)

type Server struct {
	wins          int
	losses        int
	gameID        int
	games         map[int]*battleship.Game
	previousGuess map[int]battleship.Location
}

func NewServer() *Server {
	server := &Server{
		wins:          0,
		losses:        0,
		gameID:        0,
		games:         make(map[int]*battleship.Game),
		previousGuess: make(map[int]battleship.Location),
	}
	return server
}

func (s *Server) newGame(w http.ResponseWriter, r *http.Request) {
	s.gameID++
	fmt.Fprintf(w, "{\"game_id\": %v}", s.gameID)
}

func (s *Server) turn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var turn common.Turn
	err := decoder.Decode(&turn)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Println("RAW TURN RESPONSE:", turn)

	if _, ok := s.games[turn.GameID]; !ok {
		s.games[turn.GameID] = battleship.NewGame()
		s.previousGuess[turn.GameID] = "J1"
	}

	opponentResponse := s.games[turn.GameID].HandleAttackResponse(s.previousGuess[turn.GameID], &turn)
	log.Println("RAW RESPONSE:", opponentResponse)

	guess := s.games[turn.GameID].GetRecommendedNextMove()
	s.previousGuess[turn.GameID] = guess
	log.Println("GUESS:", guess)
	responseTurn := &common.Turn{
		GameID: turn.GameID,
		Guess: common.Guess{
			Guess: string(guess),
		},
		Response: *opponentResponse,
	}
	b, err := json.Marshal(responseTurn)

	fmt.Println(s.games[turn.GameID].MyBoard)
	fmt.Println(s.games[turn.GameID].OpponentBoard)

	if turn.Response.Lost {
		log.Println("###############")
		log.Println("YOU WON!")
		log.Println("###############")
		s.wins++
		log.Printf("Record: %v wins and %v losses", s.wins, s.losses)
	}
	if opponentResponse.Lost {
		log.Println("###############")
		log.Println("YOU LOST :(")
		log.Println("###############")
		s.losses++
		log.Printf("Record: %v wins and %v losses", s.wins, s.losses)
	}
	fmt.Fprint(w, string(b))
}

func (s *Server) Run() {
	http.HandleFunc("/new_game", s.newGame)

	http.HandleFunc("/turn", s.turn)

	log.Fatal(http.ListenAndServe(":8888", nil))
}
