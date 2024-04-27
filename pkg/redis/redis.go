package redis

import (
	"be-project/pkg/config"
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/redis/go-redis/v9"
)

var storage Storage

// Storage interface that is implemented by storage providers
type Storage struct {
	db redis.UniversalClient
}

// New creates a new redis storage
func New(config config.Config) {
	// Create Universal Client
	db := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)},
		Username: config.Redis.Username,
		Password: config.Redis.Password,
		PoolSize: 10 * runtime.GOMAXPROCS(0),
	})

	// Test connection
	if err := db.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	// Create new store
	storage = Storage{
		db: db,
	}
}

func GetStorage() Storage {
	return storage
}

// Get value by key
func (s Storage) Get(key string) ([]byte, error) {
	if len(key) <= 0 {
		return nil, nil
	}
	val, err := s.db.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

// Set key with value
func (s Storage) Set(key string, val []byte, exp time.Duration) error {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}
	return s.db.Set(context.Background(), key, val, exp).Err()
}

// Delete key by key
func (s Storage) Delete(key string) error {
	if len(key) <= 0 {
		return nil
	}
	return s.db.Del(context.Background(), key).Err()
}

// Reset all keys
func (s Storage) Reset() error {
	return s.db.FlushDB(context.Background()).Err()
}

// Close the database
func (s Storage) Close() error {
	return s.db.Close()
}
