package main

import (
	"fmt"
	"game-poc/client/internal/network"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello, this is the POC game client.")
	log.SetFlags(0)

	// Something like a singleton for the game client
	client := network.GameClient
	defer client.Close()

	// listen for updates from server
	go client.ListenRead()

	// simulating game loop
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// send update to server
	client.ListenWrite(ticker.C)
}
