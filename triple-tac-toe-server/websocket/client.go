package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"gregvader/triple-tac-toe/domain"
	"gregvader/triple-tac-toe/websocket/messageType"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	User domain.User
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type messageType.MessageType `json:"type"`
	Body string                  `json:"body"`
}

type TurnMessage struct {
	Type messageType.MessageType `json:"type"`
	Body []int                   `json:"body"`
}

type StartGameMessageBody struct {
	Opponent    string `json:"opponent"`
	Number      int    `json:"number"`
	StartNumber int    `json:"startNumber"`
}

func NewWaitMessage() Message {
	return Message{Type: messageType.Wait}
}

func NewUsernameTakenMessage() Message {
	return Message{Type: messageType.UsernameTaken}
}

func NewStartGameMessage(opponent string, number int, startNumber int) Message {
	b, _ := json.Marshal(
		StartGameMessageBody{
			Opponent:    opponent,
			Number:      number,
			StartNumber: startNumber,
		})
	return Message{
		Type: messageType.StartGame,
		Body: string(b)}
}

func NewLeftMessage() Message {
	return Message{Type: messageType.Left}
}

type WebsocketMessage struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

type ClientMessage struct {
	Sender           *Client
	WebsocketMessage WebsocketMessage
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		message := WebsocketMessage{Type: messageType, Body: string(p)}
		clientMessage := ClientMessage{Sender: c, WebsocketMessage: message}
		c.Pool.Broadcast <- clientMessage
		fmt.Printf("Message Received: %+v\n", message)
	}
}
