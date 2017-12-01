package client

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/dan-v/seattlerb-battleship/api"
	"github.com/dan-v/seattlerb-battleship/battleship"
	"github.com/dan-v/seattlerb-battleship/common"
)

func RunClient(address string) {
	game := battleship.NewGame()

	client, err := api.NewAPI("http://" + address + "/")
	if err != nil {
		log.Fatalln("Failed to start game", err)
	}

	var resp *common.Response
	for {
		guess := game.GetRecommendedNextMove()
		log.Println("GUESS:", guess)
		turn, err := client.Attack(string(guess), resp)
		log.Println("RAW TURN RESPONSE:", turn)
		if err != nil {
			log.Fatalln("Failed to take turn", err)
		}

		resp = game.HandleAttackResponse(guess, turn)
		log.Println("RAW ATTACK RESPONSE:", turn)

		fmt.Println(game.MyBoard)
		fmt.Println(game.OpponentBoard)
		if resp.Lost {
			log.Println("YOU LOST :(")
			client.Attack(string(guess), resp)
			os.Exit(0)
		}
		if turn.Response.Lost {
			log.Println("YOU WON!!")
			client.Attack(string(guess), resp)
			os.Exit(0)
		}
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
	}
}
