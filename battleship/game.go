package battleship

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dan-v/seattlerb-battleship/common"
)

const (
	SEARCH_MODE = "search"
	ATTACK_MODE = "attack"
)

type Game struct {
	OpponentBoard   *Board
	opponentShips   *Ships
	MyBoard         *Board
	mode            string
	attackLocations []Location
	myBoardAttacks  map[string]bool
	turns           int
	openHits        int
}

func NewGame() *Game {
	game := &Game{
		OpponentBoard:   NewOpponentBoard(10),
		opponentShips:   NewShips(),
		MyBoard:         NewMyBoard(10),
		mode:            SEARCH_MODE,
		attackLocations: []Location{},
		myBoardAttacks:  make(map[string]bool, 0),
		turns:           0,
		openHits:        0,
	}
	game.MyBoard.SetupShips()
	return game
}

func (g *Game) GetMode() string {
	if len(g.attackLocations) == 0 {
		return SEARCH_MODE
	}
	return ATTACK_MODE
}

func (g *Game) HandleAttackResponse(guess Location, turn *common.Turn) *common.Response {
	g.turns++
	if turn.Response.Hit {
		g.openHits++

		log.Println("RESULT: HIT")
		g.OpponentBoard.MarkHit(guess)

		if turn.Response.Sunk == "" {
			g.attackLocations = append(g.attackLocations, g.OpponentBoard.GetLocationsAround(guess)...)
		} else {
			sunkShip := turn.Response.Sunk
			log.Println("SUNK SHIP:", turn.Response.Sunk)
			sunkShipSize := g.opponentShips.ShipSize(sunkShip)
			if g.openHits > sunkShipSize {
				log.Printf("Staying in Target mode. Number of hits (%v) is greater than sunk ship %s size (%v)", g.openHits, sunkShip, sunkShipSize)
				g.openHits -= sunkShipSize
			} else {
				g.openHits = 0
				g.attackLocations = []Location{}
			}
			g.opponentShips.RemoveShip(sunkShip)
		}
	} else {
		log.Println("RESULT: MISS")
		g.OpponentBoard.MarkMiss(guess)
	}
	if turn.Response.Lost {
		log.Printf("You won!")
		fmt.Println(g.MyBoard)
		fmt.Println(g.OpponentBoard)
	}
	log.Println("Attack locations:", g.attackLocations)

	response := &common.Response{
		Hit:  false,
		Lost: false,
		Sunk: "",
	}
	if _, ok := g.myBoardAttacks[turn.Guess.Guess]; ok {
		log.Printf("DUPLICATE opponent attack on %s, returning hit as %v", turn.Guess.Guess, g.myBoardAttacks[turn.Guess.Guess])
		response.Hit = g.myBoardAttacks[turn.Guess.Guess]
		return response
	}
	log.Println("OPPONENT GUESS:", turn.Guess.Guess)
	g.myBoardAttacks[turn.Guess.Guess] = false
	ship := g.MyBoard.CheckForShip(Location(turn.Guess.Guess))
	if ship != nil {
		response.Hit = true
		ship.Health--
		g.myBoardAttacks[turn.Guess.Guess] = true
		log.Println("OPPONENT RESULT: HIT", ship.Name)
		g.MyBoard.MarkHit(Location(turn.Guess.Guess))
		if ship.Health <= 0 {
			response.Sunk = ship.Name
			log.Println("OPPONENT RESULT: SUNK", ship.Name)
			g.MyBoard.shipsRemaining--
			if g.MyBoard.shipsRemaining <= 0 {
				response.Lost = true
			}
		}
	} else {
		log.Println("OPPONENT RESULT: MISS")
		g.MyBoard.MarkMiss(Location(turn.Guess.Guess))
	}
	return response
}

func (g *Game) GetRecommendedNextMove() Location {
	var guess Location
	if g.GetMode() == SEARCH_MODE {
		guess = g.GetHighestProbability()
		log.Println("MODE: SEARCH")
	} else {
		for {
			guess, g.attackLocations = g.attackLocations[len(g.attackLocations)-1], g.attackLocations[:len(g.attackLocations)-1]
			if !g.OpponentBoard.CheckLocationUnplayed(guess) {
				log.Println("CLEANING UP ALREADY PLAYED MOVE:", guess)
			} else {
				break
			}
		}
		log.Println("MODE: ATTACK")
	}
	return guess
}

func (g *Game) GetHighestProbability() Location {
	length := g.opponentShips.BiggestRemainingShip()
	probability := make([][]int, g.OpponentBoard.size)
	for i := 0; i < g.OpponentBoard.size; i++ {
		probability[i] = make([]int, g.OpponentBoard.size)

		for j := range probability[i] {
			probability[i][j] = 0
		}
	}

	for row := 0; row < g.OpponentBoard.size; row++ {
		for col := 0; col+length <= g.OpponentBoard.size; col++ {
			blocker := false
			for i := col; i < col+length; i++ {
				if g.OpponentBoard.layout[row][i] != MARK_EMPTY {
					blocker = true
				}
			}
			if !blocker {
				for i := col; i < col+length; i++ {
					probability[row][i]++
				}
			}
		}
	}

	for col := 0; col < g.OpponentBoard.size; col++ {
		for row := 0; row+length <= g.OpponentBoard.size; row++ {
			blocker := false
			for i := row; i < row+length; i++ {
				if g.OpponentBoard.layout[i][col] != MARK_EMPTY {
					blocker = true
				}
			}
			if !blocker {
				for i := row; i < row+length; i++ {
					probability[i][col]++
				}
			}
		}
	}

	output := " "
	for i := 1; i <= g.OpponentBoard.size; i++ {
		output += fmt.Sprint(i, " ")
	}
	output += fmt.Sprintln("")
	for i, row := range probability {
		output += fmt.Sprintf("%c ", 'A'+i)
		for _, irow := range row {
			if irow >= 10 {
				output += fmt.Sprintf("%v", irow)
			} else {
				output += fmt.Sprintf("%v ", irow)
			}
		}
		output += fmt.Sprintln("")
	}

	highest := 0
	var highRow, highCol int
	for i := 0; i < g.OpponentBoard.size; i++ {
		for j := 0; j < g.OpponentBoard.size; j++ {
			if probability[i][j] > highest {
				highest = probability[i][j]
				highRow = i
				highCol = j
			}
		}
	}
	highestProbability := cordinatesToLocation(highRow+1, highCol+1)
	return highestProbability
}

func (g *Game) GetRandomUnusedGuess() Location {
	for {
		row := randomInt(g.OpponentBoard.size)
		col := randomInt(g.OpponentBoard.size)
		guess := cordinatesToLocation(row, col)
		if g.OpponentBoard.CheckLocationUnplayed(guess) {
			return guess
		}
	}
}

func randomInt(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max) + 1
}
