package entities

type Auth struct {
	ID       int `gorm:"primaryKey"`
	Username string
	Password string
	Token    *string
}
