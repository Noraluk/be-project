package entities

import "time"

type ChatUser struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Username  string    `json:"username"`
	IsLoggin  bool      `json:"is_loggin"`
	UpdatedAt time.Time `json:"updated_at"`
}
