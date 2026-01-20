package redisclient

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/redis/go-redis/v9"
)

const TaskQueue = "task_queue"

type Client struct {
	rdb *redis.Client
}

func New() *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &Client{rdb: rdb}
}

func (c *Client) PushTask(ctx context.Context, task any) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return c.rdb.LPush(ctx, TaskQueue, data).Err()
}

func (c *Client) PopTask(ctx context.Context) ([]byte, error) {
	res, err := c.rdb.BRPop(ctx, 0, TaskQueue).Result()
	if err != nil {
		return nil, err
	}
	if len(res) != 2 {
		return nil, errors.New("invalid BRPOP response")
	}
	return []byte(res[1]), nil
}
