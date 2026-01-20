package main

import (
	"context"
	"fmt"

	"Jobqueue/internal/redisclient"
)

type Task struct {
	Type string
}

func main() {
	ctx := context.Background()
	redis := redisclient.New()

	task := Task{Type: "test"}

	err := redis.PushTask(ctx, task)
	if err != nil {
		panic(err)
	}

	data, err := redis.PopTask(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
