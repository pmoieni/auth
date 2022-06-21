package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisDTO struct {
	Key   string
	Value string
	Exp   time.Duration
}

var RefreshTokenDB *redis.Client
var ClientIDDB *redis.Client
var PasswordResetTokenDB *redis.Client

func CreateRefreshTokenDatabase(addr, password string) (err error) {
	RefreshTokenDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	_, err = RefreshTokenDB.Ping(context.Background()).Result()
	return
}

func CreateClientIDDatabase(addr, password string) (err error) {
	ClientIDDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1,
	})
	_, err = ClientIDDB.Ping(context.Background()).Result()
	return
}

func CreatePasswordResetTokenDatabase(addr, password string) (err error) {
	PasswordResetTokenDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       2,
	})
	_, err = PasswordResetTokenDB.Ping(context.Background()).Result()
	return
}

func (r *RedisDTO) Set(db *redis.Client) error {
	_, err := db.Set(context.Background(), r.Key, r.Value, r.Exp).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisDTO) Get(db *redis.Client) (string, error) {
	value, err := db.Get(context.Background(), r.Key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

func (r *RedisDTO) Delete(db *redis.Client) error {
	_, err := db.Del(context.Background(), r.Key).Result()
	if err != nil {
		return err
	}
	return nil
}
