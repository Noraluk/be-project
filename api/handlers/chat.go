package handlers

import (
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	Username string `json:"username"`
	Target   string `json:"target"`
	Message  string `json:"message"`
}

type Client struct {
	Username string
	Conn     *websocket.Conn
}

type ChatHandler interface {
	CreateConnection()
	Broadcast(c *websocket.Conn)
}

type chatHandler struct {
	clients    map[string]*websocket.Conn
	register   chan Client
	broadcast  chan Message
	unregister chan Client
	users      []string
}

func NewChatHandler() ChatHandler {
	return &chatHandler{
		clients:    make(map[string]*websocket.Conn),
		register:   make(chan Client),
		broadcast:  make(chan Message),
		unregister: make(chan Client),
		users:      make([]string, 0),
	}
}

func (h chatHandler) CreateConnection() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.Username] = client.Conn
			h.users = append(h.users, client.Username)
			h.notifyClients()

		case message := <-h.broadcast:
			source := h.clients[message.Username]
			destination, ok := h.clients[message.Target]

			if err := source.Conn.WriteJSON(message); err != nil {
				log.Println("write:", err)
			}
			if ok {
				if err := destination.Conn.WriteJSON(message); err != nil {
					log.Println("write:", err)
				}
			}

		case client := <-h.unregister:
			delete(h.clients, client.Username)
			var i int
			for j, user := range h.users {
				if user == client.Username {
					i = j
				}
			}
			h.users = append(h.users[:i], h.users[i+1:]...)
			h.notifyClients()
		}
	}
}

func (h chatHandler) Broadcast(c *websocket.Conn) {
	ct := Client{
		Username: c.Locals("username").(string),
		Conn:     c,
	}

	defer func() {
		h.unregister <- ct
		c.Close()
	}()

	h.register <- ct

	fmt.Println("foo")

	for {
		mt, m, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", m)

		if mt == websocket.TextMessage {
			h.broadcast <- Message{
				Username: c.Locals("username").(string),
				Target:   c.Locals("target").(string),
				Message:  string(m),
			}
		}
	}
}

func (h chatHandler) notifyClients() {
	fmt.Println(h.clients)
	for username, client := range h.clients {
		clientIDs := []string{}
		for _, un := range h.users {
			if un != username {
				clientIDs = append(clientIDs, un)
			}
		}
		err := client.WriteJSON(clientIDs)
		if err != nil {
			log.Println("notify error:", err)
		}
	}
}
