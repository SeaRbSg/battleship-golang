package main

import (
	"flag"
	"log"

	"github.com/dan-v/seattlerb-battleship/client"
	"github.com/dan-v/seattlerb-battleship/server"
)

func main() {
	clientOrServer := flag.String("type", "", "client/server")
	address := flag.String("address", "127.0.0.1:8888", "server address to connect to")
	flag.Parse()

	if *clientOrServer != "client" && *clientOrServer != "server" {
		log.Fatalln("Type must be 'client' or 'server'")
	}

	if *clientOrServer == "client" {
		client.RunClient(*address)
	} else {
		server := server.NewServer()
		server.Run()
	}
}
