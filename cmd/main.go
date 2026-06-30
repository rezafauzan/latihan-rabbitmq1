package main

import (
	"log"
	"math/rand/v2"
	"rezafauzan/latihan-rabbitmq/internal/lib"
	"time"
)

func main() {
	rabbit, err := lib.NewRabbit()
	if err != nil {
		log.Fatal("Rabbit Fail To Establish")
	}
	
	for true {
		err = lib.ProduceMessage(rabbit.RMQ, "dev.test", "ini itu")
		if err != nil {
			log.Fatal("Rabbit Fail To Produce Message", err)
		}
	
		log.Printf("OK")
		time.Sleep(time.Duration(rand.Float64()) * time.Second)
	}

}
