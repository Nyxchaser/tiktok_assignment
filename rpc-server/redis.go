package main

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	cli *redis.Client
}

func (c *RedisClient) InitClient(ctx context.Context, address, password string) error {
	r := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	// test connection
	if err := r.Ping(ctx).Err(); err != nil {
		return err
	}

	c.cli = r
	return nil
}

func (c *RedisClient) SaveMessage(ctx context.Context, roomID string, message *Message) error {
	// Marshal the Go struct into JSON bytes
	text, err := json.Marshal(message)
	if err != nil {
		return err
	}

	member := &redis.Z{
		Score:  float64(message.Timestamp), // The sort key
		Member: text,                       // Data
	}

	_, err = c.cli.ZAdd(ctx, roomID, *member).Result()
	if err != nil {
		return err
	}

	return nil
}
