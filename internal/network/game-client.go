package network

import (
	"fmt"
	"game-poc/client/internal/configurations"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	connection *websocket.Conn
	done       chan struct{}
}

var GameClient = initClient()

func initClient() *Client {
	c := configurations.Config
	url := url.URL{Scheme: "ws", Host: c.GameServer.Host, Path: c.GameServer.Path}
	return &Client{
		connection: getWebsocketConnection(url),
		done:       make(chan struct{}),
	}
}

func getWebsocketConnection(url url.URL) *websocket.Conn {
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatalf("Could not open websocket to '%s'", url.String())
	}
	return c
}

func (c *Client) Close() {
	c.connection.Close()
}

func (c *Client) ListenRead() {
	defer close(c.done)
	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			log.Println("Game Client Websocket Closed: ", err)
			return
		}
		log.Printf("Received From Server: %s", message)
	}
}

func (c *Client) ListenWrite(ch <-chan time.Time) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-c.done:
			return
		case m := <-ch:
			message := []byte(m.String())
			c.handleWrite(message)
		case <-interrupt:
			c.handleInterrupt()
		}
	}
}

func (c *Client) handleWrite(message []byte) {
	if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println("Websocket Write Error: ", err)
		if ce, ok := err.(*websocket.CloseError); ok {
			switch ce.Code {
			case websocket.CloseNormalClosure:
				log.Println("Websocket Closed By Client - exiting")
				return
			}
		}
	}
}

func (c *Client) handleInterrupt() {
	log.Println("Interupt Received, will attempt to close connection to server")
	err := c.connection.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("Websocket CloseMessage Error: ", err)
		return
	}
	select {
	case <-c.done:
		fmt.Println("Done closing websocket")
	case <-time.After(time.Second):
		fmt.Println("Received timeout on attempt to close websocket")
	}
}
