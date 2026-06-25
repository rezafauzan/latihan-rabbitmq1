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

	err = lib.ProduceMessage(rabbit.RMQ, "dev.test", "ini itu")
	if err != nil {
		log.Fatal("Rabbit Fail To Produce Message", err)
	}

	err = lib.ProduceMessage(rabbit.RMQ, "dev.test", "ini itu")
	if err != nil {
		log.Fatal("Rabbit Fail To Produce Message", err)
	}
	log.Printf("OK")
}
