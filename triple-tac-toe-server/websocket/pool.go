package websocket

import (
	"fmt"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Room       *Room
	Broadcast  chan ClientMessage
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Room:       NewRoom(),
		Broadcast:  make(chan ClientMessage),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Room.Add(client)
		case client := <-pool.Unregister:
			pool.Room.Remove(client)
			fmt.Println("WS> Unegistered User:" + fmt.Sprint(client.User.Username))
		case clientMessage := <-pool.Broadcast:
			//var userIds []uint64
			pool.Room.SendMessage(clientMessage)

			/*client := &http.Client{}
			req, _ := http.NewRequest("GET", globals.GetUsersMicroserviceUrl()+"/api/users/notify", nil)
			authInfo, _ := tokens.CreateAuthInfo(domain.User{Id: clientMessage.Client.UserId})
			req.Header.Set("Authorization", "Bearer "+authInfo.Token)
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}

			json.NewDecoder(resp.Body).Decode(&userIds)
			*/
			/*
				for _, id := range userIds {
					if client := pool.Clients[id]; client != nil {
						if err := client.Conn.WriteJSON(clientMessage.Message); err != nil {
							fmt.Println(err)
							return
						}
					}
				}*/
		}
	}
}
