// dummyengine/pkg/rabbitmqQueue/consumer.go
package rabbitmqQueue

import (
	"dummyengine/pkg/orderbook"
	"dummyengine/pkg/pricelevel"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// RabbitMQConsumer handles message consumption
type RabbitMQQueue struct {
	Conn      *amqp.Connection
	Ch        *amqp.Channel
	Queue     string
	URL       string
	OrderBook *orderbook.OrderBook // OrderBook instance
}

// OrderMessage represents the structure received from RabbitMQ
type OrderMessage struct {
	Type  string           `json:"type"` // "buy" or "sell"
	Order pricelevel.Order `json:"order"`
}

// NewRabbitMQConsumer initializes a new consumer
func NewRabbitMQQueue(url, queue, tradeQueue, priceQueue, orderBookQueue string) *RabbitMQQueue {
	return &RabbitMQQueue{
		Queue:     queue,
		URL:       url,
		OrderBook: nil, // Initialize to nil, will be set after connection
	}
}

// Connect establishes the RabbitMQ connection
func (c *RabbitMQQueue) Connect() {
	var err error

	// Attempt connection
	c.Conn, err = amqp.Dial(c.URL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}

	c.Ch, err = c.Conn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open a channel: %v", err)
	}

	log.Printf("✅ Connected to RabbitMQ: %s (Queue: %s)", c.URL, c.Queue)

	// Auto-reconnect on failure
	c.Conn.NotifyClose(make(chan *amqp.Error))
	c.reconnect()

	// Start consuming
	c.Consume()
}

// Reconnect handles RabbitMQ reconnection logic
func (c *RabbitMQQueue) reconnect() {
	for {
		log.Printf("🔄 Reconnecting to RabbitMQ...")
		time.Sleep(5 * time.Second)

		var err error
		c.Conn, err = amqp.Dial(c.URL)
		if err == nil {
			c.Ch, err = c.Conn.Channel()
			if err == nil {
				log.Println("✅ Reconnected to RabbitMQ")
				c.OrderBook = orderbook.NewOrderBook(c.Ch, "trade_queue", "price_queue", "orderBook_queue")
				c.Consume()
				return
			}
		}
	}
}

// Consume listens to the queue
func (c *RabbitMQQueue) Consume() {
	msgs, err := c.Ch.Consume(
		c.Queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Failed to register a consumer: %v", err)
	}

	log.Printf("📥 Consuming messages from queue: %s", c.Queue)

	for msg := range msgs {
		c.processMessage(msg)
	}
}

// processMessage handles incoming messages and updates OrderBook
func (c *RabbitMQQueue) processMessage(msg amqp.Delivery) {
	var orderMsg OrderMessage
	if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
		log.Printf("❌ Failed to parse message: %v", err)
		msg.Nack(false, false) // Reject message
		return
	}

	log.Printf("📦 Order Processed: %+v", orderMsg.Order)

	// Process order based on type
	switch orderMsg.Type {
	case "buy":
		c.OrderBook.AddBuyOrder(&orderMsg.Order)
	case "sell":
		c.OrderBook.AddSellOrder(&orderMsg.Order)
	default:
		log.Printf("⚠️ Unknown order type: %s", orderMsg.Type)
		msg.Nack(false, false)
		return
	}

	msg.Ack(false) // Acknowledge message
}
