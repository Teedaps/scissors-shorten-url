package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"os"
)

var ctx = context.Background()

func createclient(dbNo int) *redis.client {
	rdb := redis.Newclient(&redis.options{
		Addr:     os.Getenv("DB_ADDR"),
		password: os.Getenv("DB_PASS"),
		DB:       dbNo,
	})
	return rdb
}
