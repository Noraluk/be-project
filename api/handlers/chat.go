package handlers

import (
	"be-project/api/entities"
	"be-project/pkg/base"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
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
	repository base.BaseRepository[any]
}

func NewChatHandler(repository base.BaseRepository[any]) ChatHandler {
	return &chatHandler{
		clients:    make(map[string]*websocket.Conn),
		register:   make(chan Client),
		broadcast:  make(chan Message),
		unregister: make(chan Client),
		users:      make([]string, 0),
		repository: repository,
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
			sender := h.clients[message.Sender]
			recipient, ok := h.clients[message.Recipient]

			chat := entities.Chat{
				Sender:    message.Sender,
				Recipient: message.Recipient,
				Message:   message.Message,
				CreatedAt: time.Now(),
			}
			if err := h.repository.Create(&chat).Error(); err != nil {
				log.Println("create chat failed, error: ", err)
			}

			if err := sender.Conn.WriteJSON(message); err != nil {
				log.Println("write:", err)
			}
			if ok {
				if err := recipient.Conn.WriteJSON(message); err != nil {
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
		Username: c.Locals("sender").(string),
		Conn:     c,
	}

	defer func() {
		h.unregister <- ct
		c.Close()
	}()

	h.register <- ct

	recipient := c.Locals("recipient").(string)
	if len(recipient) > 0 {
		var chats []entities.Chat
		if err := h.repository.Where("(sender = '1' and recipient = '2') or (sender = '2' and recipient = '1')").Order("id desc").Limit(50).Find(&chats).Error(); err != nil {
			log.Println("find chats failed, error: ", err)
		}

		for i := len(chats) - 1; i >= 0; i-- {
			err := c.WriteJSON(Message{Sender: chats[i].Sender, Recipient: chats[i].Recipient, Message: chats[i].Message})
			if err != nil {
				log.Println("write failed, error: ", err)
			}
		}
	}

	for {
		mt, m, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", m)

		if mt == websocket.TextMessage {
			h.broadcast <- Message{
				Sender:    c.Locals("sender").(string),
				Recipient: recipient,
				Message:   string(m),
			}
		}
	}
}

func (h chatHandler) notifyClients() {
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
