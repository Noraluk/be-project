package dtos

type ChatUser struct {
	Username    string `json:"username"`
	UnreadCount int    `json:"unread_count"`
}
