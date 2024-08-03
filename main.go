package main

import (
	"log"

	"github.com/risersh/controller/conf"
	"github.com/risersh/controller/monitoring"
	"github.com/risersh/controller/rabbitmq"
)

func init() {
	conf.Init()
	monitoring.Setup()
}

func main() {
	go rabbitmq.Setup()

	ch, err := rabbitmq.StartConsuming[rabbitmq.Message[rabbitmq.MessageType]]()
	if err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}
	go func() {
		for msg := range ch {
			log.Println(msg)
		}
	}()
}
