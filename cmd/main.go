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

	err = lib.ProduceMessage(rabbit.RMQ, "dev.test", "ini itu pokoknya kirim")
	if err != nil {
		log.Fatal("Rabbit Fail To Produce Message", err)
	}

	err = lib.ConsumeMessage(rabbit.RMQ, "dev.test", "test.consumer")
	if err != nil {
		log.Fatal("Rabbit Fail To Consume Message", err)
	}

	err = lib.ProduceMessage(rabbit.RMQ, "dev.test", "ini itu pokoknya kirim")
	if err != nil {
		log.Fatal("Rabbit Fail To Produce Message", err)
	}

}
