// dummyengine/cmd/main.go
package main

import (
	"dummyengine/pkg/pricelevel"
	"dummyengine/pkg/rabbitmqQueue"
)

// OrderMessage represents the structure sent to RabbitMQ
type OrderMessage struct {
	Type     string           `json:"type"` // "buy" or "sell"
	Order    pricelevel.Order `json:"order"`
}

func main() {
	// Initialize RabbitMQ Publisher
	consumer := rabbitmqQueue.NewRabbitMQQueue("amqp://guest:guest@localhost:5672/", "order_queue", "trade_queue", "price_queue", "orderBook_queue")
	consumer.Connect()

	select {}
}
