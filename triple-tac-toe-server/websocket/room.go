package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"gregvader/triple-tac-toe/websocket/messageType"
	"math/rand"
	"sync"
)

type Room struct {
	clients     map[string]*Client
	waitingRoom map[string]*Client
	opponent    sync.Map
	mutex       *sync.RWMutex
}

func NewRoom() *Room {
	return &Room{
		clients:     make(map[string]*Client),
		waitingRoom: make(map[string]*Client),
		opponent:    sync.Map{}, //map[string]*Client
		mutex:       &sync.RWMutex{},
	}
}

func (r *Room) Add(client *Client) {
	r.mutex.Lock()
	if r.clients[client.User.Username] != nil {
		message := NewUsernameTakenMessage()
		client.Conn.WriteJSON(message)
	} else {
		r.clients[client.User.Username] = client
		fmt.Println("WS> Registered User:" + fmt.Sprint(client.User.Username))

		opponent, err := r.getFirstWaitingOpponent()
		if err != nil {
			r.waitingRoom[client.User.Username] = client
			message := NewWaitMessage()
			client.Conn.WriteJSON(message)
		} else {
			//Pait the two players and start the game
			opponentClient := r.clients[opponent]
			delete(r.waitingRoom, opponent)

			r.opponent.Store(client.User.Username, opponentClient)
			r.opponent.Store(opponent, client)

			player1, player2 := r.getPlayerNumbers()

			message1 := NewStartGameMessage(opponent, player1, 1)
			message2 := NewStartGameMessage(client.User.Username, player2, 1)
			client.Conn.WriteJSON(message1)
			opponentClient.Conn.WriteJSON(message2)
		}
	}
	r.mutex.Unlock()
}

func (r *Room) getFirstWaitingOpponent() (string, error) {
	for k := range r.waitingRoom {
		return k, nil
	}
	return "", errors.New("waiting room is empty")
}

func (r *Room) getPlayerNumbers() (int, int) {
	first := rand.Intn(2)
	if first == 0 {
		return 0, 1
	} else {
		return 1, 0
	}
}

func (r *Room) IsUsernameTaken(username string) bool {
	return r.clients[username] != nil
}

func (r *Room) Remove(client *Client) {
	r.mutex.Lock()
	delete(r.clients, client.User.Username)
	delete(r.waitingRoom, client.User.Username)
	opponent, present := r.opponent.LoadAndDelete(client.User.Username)
	if present {
		opponentClient := opponent.(*Client)
		r.opponent.Delete(opponentClient.User.Username)
		delete(r.waitingRoom, opponentClient.User.Username)
		delete(r.clients, opponentClient.User.Username)

		message := NewLeftMessage()
		opponentClient.Conn.WriteJSON(message)
		//opponentClient.Conn.Close()
	}
	r.mutex.Unlock()
}

func (r *Room) SendMessage(clientMessage ClientMessage) {
	var message Message
	jsonErr := json.Unmarshal([]byte(clientMessage.WebsocketMessage.Body), &message)

	if jsonErr == nil {
		if message.Type == messageType.KeepAlive {
			clientMessage.Sender.Conn.WriteJSON(Message{Type: messageType.Ok})
			return
		}
	}

	opponent, present := r.opponent.Load(clientMessage.Sender.User.Username)
	if present {
		opponentClient := opponent.(*Client)
		var message = TurnMessage{}
		jsonErr := json.Unmarshal([]byte(clientMessage.WebsocketMessage.Body), &message)
		if jsonErr != nil {
			fmt.Println(jsonErr.Error())
		}
		fmt.Println(message.Type)
		fmt.Println(message.Body)
		arrayString, _ := json.Marshal(message.Body)
		opponentClient.Conn.WriteJSON(Message{Type: message.Type, Body: string(arrayString)})
	}
}
