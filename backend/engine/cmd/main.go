package main

import (
	"dummyengine/pkg/config"
	"dummyengine/pkg/rabbitmqQueue"
	"log"
	"dummyengine/pkg/logger"
)

func main() {
	logger.InitLogger()
	config.LoadConfig()
	
	rabbitMQURL := config.AppConfig.RabbitMQ.URL
	if rabbitMQURL == "" {
		log.Fatal("RabbitMQ URL is not set in the environment variables")
	}

	// Initialize RabbitMQ Publisher
	consumer := rabbitmqQueue.NewRabbitMQQueue(
		rabbitMQURL,
		"order_queue",
		"trade_queue",
		"price_queue",
		"orderBook_queue",
	)
	consumer.Connect()

	select {}
}
