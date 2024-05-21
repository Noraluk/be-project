package handlers

import (
	"be-project/api/dtos"
	"be-project/api/entities"
	"be-project/pkg/base"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type RequestType int

const (
	ChatHistory RequestType = 1
	Chat        RequestType = 2
)

type Message struct {
	Sender      string      `json:"sender"`
	Recipient   string      `json:"recipient"`
	Message     string      `json:"message"`
	RequestType RequestType `json:"request_type"`
	Unread      bool        `json:"unread"`
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
	repository base.BaseRepository[any]
}

func NewChatHandler(repository base.BaseRepository[any]) ChatHandler {
	return &chatHandler{
		clients:    make(map[string]*websocket.Conn),
		register:   make(chan Client),
		broadcast:  make(chan Message),
		unregister: make(chan Client),
		repository: repository,
	}
}

func (h chatHandler) CreateConnection() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.Username] = client.Conn

			var chatUser entities.ChatUser
			err := h.repository.First(&chatUser, "username = ?", client.Username).Error()
			if err == nil {
				err = h.repository.Model(&entities.ChatUser{}).Where("username = ?", client.Username).Update("is_loggin", true).Error()
				if err != nil {
					log.Println("update is loggin failed, eror: ", err)
				}
			} else {
				cu := entities.ChatUser{
					Username:  client.Username,
					IsLoggin:  true,
					UpdatedAt: time.Now(),
				}

				err = h.repository.Create(&cu).Error()
				if err != nil {
					log.Println("create chat user failed, error: ", err)
				}
			}

			h.notifyClients()

		case message := <-h.broadcast:
			sender := h.clients[message.Sender]
			recipient, ok := h.clients[message.Recipient]

			chat := entities.Chat{
				Sender:    message.Sender,
				Recipient: message.Recipient,
				Message:   message.Message,
				Unread:    message.Unread,
				CreatedAt: time.Now(),
			}
			if err := h.repository.Create(&chat).Error(); err != nil {
				log.Println("create chat failed, error: ", err)
			}

			if err := sender.Conn.WriteJSON(message); err != nil {
				log.Println("write:", err)
			}
			if ok {
				h.online(recipient, message.Recipient)

				if err := recipient.Conn.WriteJSON(message); err != nil {
					log.Println("write:", err)
				}

			}

		case client := <-h.unregister:
			var chatUser entities.ChatUser
			err := h.repository.First(&chatUser, "username = ?", client.Username).Error()
			if err == nil {
				err = h.repository.Model(&chatUser).Update("is_loggin", false).Error()
				if err != nil {
					log.Println("update is loggin failed, eror: ", err)
				}
			}

			delete(h.clients, client.Username)
			h.notifyClients()
		}
	}
}

type Request struct {
	Sender string `json:"sender"`
}

func (h chatHandler) Broadcast(c *websocket.Conn) {
	var req Request
	err := c.ReadJSON(&req)
	if err != nil {
		log.Println("read json failed, error: ", err)
	}

	ct := Client{
		Username: req.Sender,
		Conn:     c,
	}

	defer func() {
		h.unregister <- ct
		c.Close()
	}()

	h.register <- ct

	for {
		var msg Message
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %v", msg)

		switch msg.RequestType {
		case ChatHistory:
			var chats []entities.Chat
			if err := h.repository.Where("(sender = ? and recipient = ?) or (sender = ? and recipient = ?)", msg.Sender, msg.Recipient, msg.Recipient, msg.Sender).Order("id desc").Limit(50).Find(&chats).Error(); err != nil {
				log.Println("find chats failed, error: ", err)
			}

			for i := len(chats) - 1; i >= 0; i-- {
				err := c.WriteJSON(Message{Sender: chats[i].Sender, Recipient: chats[i].Recipient, Message: chats[i].Message, Unread: chats[i].Unread})
				if err != nil {
					log.Println("write failed, error: ", err)
				}
			}
			h.markMessagesAsRead(msg.Sender, msg.Recipient)
			h.online(c, msg.Sender)
		case Chat:
			msg.Unread = true
			h.broadcast <- msg
		}

	}
}

func (h chatHandler) notifyClients() {
	for username, client := range h.clients {
		h.online(client, username)
	}
}

func (h chatHandler) online(client *websocket.Conn, username string) {
	var chatUsers []dtos.ChatUser
	err := h.repository.Table("chat_users cu").
		Select("cu.username, count(case when unread = true then 1 end) as unread_count").
		Joins("left join chats c on c.sender = cu.username").
		Where("cu.username != ? and is_loggin = true", username).
		Group("cu.username").Find(&chatUsers).Error()
	if err != nil {
		log.Println("find chat users failed, error: ", err)
	}

	err = client.WriteJSON(chatUsers)
	if err != nil {
		log.Println("notify error:", err)
	}
}

func (h chatHandler) markMessagesAsRead(sender, recipient string) {
	err := h.repository.Model(&entities.Chat{}).Where("recipient = ? AND sender = ? AND unread = ?", sender, recipient, true).Update("unread", false).Error()
	if err != nil {
		log.Printf("error marking messages as read: %v", err)
	}
}
