package main

import (
	"Jobqueue/internal/redisclient"
	"log"
)

func main() {
	_ = redisclient.New()
	log.Println("Worker started")
}
