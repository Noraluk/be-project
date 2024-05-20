package dtos

type ChatUser struct {
	Sender      string `json:"sender"`
	Recipient   string `json:"recipient"`
	UnreadCount int    `json:"unread_count"`
}
