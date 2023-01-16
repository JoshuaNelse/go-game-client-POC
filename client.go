package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	fmt.Println("Hello, this is the POC game client.")
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	url := url.URL{Scheme: "ws", Host: *addr, Path: "/game"}
	log.Printf("Connecting to %s", url.String())

	connection, _, err := websocket.DefaultDialer.Dial(url.String(), nil)

	if err != nil {
		log.Fatal("Websocket Dial Error: ", err)
	}

	defer connection.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := connection.ReadMessage()
			if err != nil {
				log.Println("Websocket Closed: ", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			message := []byte(t.String())
			err := connection.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Websocket Write Error: ", err)
				if ce, ok := err.(*websocket.CloseError); ok {
					switch ce.Code {
					case websocket.CloseNormalClosure:
						log.Println("Websocket Closed By Client - exiting")
						return
					}
				}

			}
		case <-interrupt:
			log.Println("Interupt")

			//Cleanup connection with close message
			//waiting w/timeout for the server to close the connection
			err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Websocket CloseMessage Error: ", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
