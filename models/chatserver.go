package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	database "v/pkg"

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

		switch message.Channel {
		case "general":
			{
				typePayload := &TextMessage{}
				jsonPayload, err := json.Marshal(message.Payload)
				if err != nil {
					fmt.Println("Failed to marshal payload:", err)
					continue
				}

				err = json.Unmarshal(jsonPayload, typePayload)
				if err != nil {
					fmt.Println("Failed to unmarshal payload:", err)
					continue
				}

				if added := database.InsertMessageIntoChannel("general", message.Sender, typePayload.Data); added {
					fmt.Println("New message added to table.")
				}

				if typePayload.Data == "allmessage" {
					messages := database.GetAllMessagesOfChannel("general")
					for _, value := range messages {
						Payload := TextMessage{
							Data: value,
						}

						MainMessage := &Message{
							Sender:  "server",
							Type:    "text",
							Channel: "nil",
							Payload: Payload,
						}
						conn.WriteJSON(MainMessage)
					}
				}
			}
		case "hindi":
			{
				database.InsertMessageIntoChannel("hindi", message.Sender, message.Payload.(TextMessage).Data)
			}
		case "english":
			{
				database.InsertMessageIntoChannel("english", message.Sender, message.Payload.(TextMessage).Data)
			}
		case "bakchodi":
			{
				database.InsertMessageIntoChannel("bakchodi", message.Sender, message.Payload.(TextMessage).Data)
			}
		}
	}
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		Users:    make(map[*websocket.Conn]bool),
		Messages: make(chan Message),
	}
}
