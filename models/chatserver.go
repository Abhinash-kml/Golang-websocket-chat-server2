package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatServer struct {
	Users    map[*websocket.Conn]bool
	Messages chan Message
	mutex    sync.Mutex
}

func (c *ChatServer) AddUser(connection *websocket.Conn) {
	c.mutex.Lock()
	c.Users[connection] = true
	c.mutex.Unlock()
}

func (c *ChatServer) RemoveUser(connection *websocket.Conn) {
	c.mutex.Lock()
	delete(c.Users, connection)
	c.mutex.Unlock()
}

func (c *ChatServer) HandleMessages() {
	for {
		message := <-c.Messages

		c.mutex.Lock()
		for conn := range c.Users {
			err := conn.WriteJSON(message)
			if err != nil {
				fmt.Println(err)
				conn.Close()
				delete(c.Users, conn)
			}
		}
		c.mutex.Unlock()
	}
}

func (c *ChatServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode("Connection upgrade error")
		return
	}
	defer conn.Close()

	// Add user
	c.AddUser(conn)

	// Listen for incoming messages and send it to broadcast channel
	for {
		message := &Message{}
		err := conn.ReadJSON(message)
		if err != nil {
			c.RemoveUser(conn)
			break
		}

		c.Messages <- *message
	}
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		Users:    make(map[*websocket.Conn]bool),
		Messages: make(chan Message),
	}
}
