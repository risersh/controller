package controller

import (
	"log"

	"github.com/risersh/controller/rabbitmq"
)

func main() {
	go rabbitmq.Setup()

	ch, err := rabbitmq.StartConsuming[rabbitmq.Message[rabbitmq.MessageTypeNewDeployment]]()
	if err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}
	go func() {
		for msg := range ch {
			log.Println(msg)
		}
	}()
}
