package entities

import "time"

type Chat struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Message   string    `json:"message"`
	Unread    bool      `json:"unread"`
	CreatedAt time.Time `json:"created_at"`
}
