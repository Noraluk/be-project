package session

import (
	"be-project/pkg/redis"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
)

func New() *session.Store {
	storage := redis.GetStorage()

	store := session.New(session.Config{
		Storage:    storage,
		Expiration: 5 * time.Minute,
		KeyLookup:  "cookie:myapp_session",
	})

	return store
}
