package main

import (
	"log"
	"rezafauzan/latihan-rabbitmq/internal/lib"
)

func main() {
	rabbit, err := lib.NewRabbit()
	if err != nil {
		log.Fatal("Rabbit Fail To Establish")
	}
	
	err = lib.ConsumeMessage(rabbit.RMQ, "dev.test", "test.consumer")
	if err != nil {
		log.Fatal("Rabbit Fail To Consume Message", err)
	}
}
