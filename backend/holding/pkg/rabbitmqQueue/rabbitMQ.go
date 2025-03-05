package rabbitmqQueue

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"holding/pkg/logger"
	"log"
	"math/big"
	"time"
	"holding/pkg/common"
)

// RabbitMQConsumer handles message consumption
type RabbitMQConsumer struct {
	Conn  *amqp.Connection
	Ch    *amqp.Channel
	Queue string
	URL   string
}

// EventMessage represents the structure received from RabbitMQ
type EventMessage struct {
	Task    string `json:"task"`              // buyOrderPlaced or sellOrderPlaced
	UserID  string `json:"user_id,omitempty"` // UserID
	OrderID string `json:"orderId,omitempty"` // OrderID
	Amount  int    `json:"amount,omitempty"`  // amount
}

// NewRabbitMQConsumer initializes a new consumer with retry mechanism
func NewRabbitMQConsumer(url, queue string) (*RabbitMQConsumer, error) {
	consumer := &RabbitMQConsumer{
		Queue: queue,
		URL:   url,
	}

	if err := consumer.Connect(); err != nil {
		return nil, err
	}

	return consumer, nil
}

// Connect establishes the RabbitMQ connection and retries on failure
func (c *RabbitMQConsumer) Connect() error {
	var err error
	for attempts := 0; attempts < 5; attempts++ {
		c.Conn, err = amqp.Dial(c.URL)
		if err == nil {
			c.Ch, err = c.Conn.Channel()
			if err == nil {
				log.Printf("✅ Connected to RabbitMQ (Queue: %s)\n", c.Queue)
				return nil
			}
			c.Conn.Close()
		}
		log.Println("❌ Failed to connect to RabbitMQ, retrying in 5s...")
		time.Sleep(5 * time.Second)
	}

	return errors.New("unable to connect to RabbitMQ after multiple attempts")
}

// Consume listens to the queue and processes messages, with auto-reconnect
func (c *RabbitMQConsumer) Consume() {
	for {
		if c.Conn == nil || c.Conn.IsClosed() {
			log.Println("⚠️ RabbitMQ connection lost, attempting to reconnect...")
			if err := c.Connect(); err != nil {
				log.Println("❌ Reconnection failed, retrying in 5s...")
				time.Sleep(5 * time.Second)
				continue
			}
		}

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
			log.Printf("❌ Consumer failed: %v, retrying in 5s...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("📥 Consuming messages from queue: %s\n", c.Queue)

		for msg := range msgs {
			c.processMessage(msg)
		}

		log.Println("⚠️ Message channel closed, reconnecting in 5s...")
		time.Sleep(5 * time.Second)
	}
}

// processMessage processes incoming messages
func (c *RabbitMQConsumer) processMessage(msg amqp.Delivery) {
	var orderMsg EventMessage
	if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
		log.Printf("❌ Failed to parse message: %v\n", err)
		msg.Nack(false, false)
		return
	}

	UserID := new(big.Int)
	if err := UserID.UnmarshalText([]byte(orderMsg.UserID)); err != nil {
		log.Printf("❌ Failed to convert orderMsg.UserID to big.Int: %v\n", err)
		msg.Nack(false, false)
		return
	}

	OrderID := new(big.Int)
	if err := OrderID.UnmarshalText([]byte(orderMsg.OrderID)); err != nil {
		log.Printf("❌ Failed to convert orderMsg.OrderID to big.Int: %v\n", err)
		msg.Nack(false, false)
		return
	}

	switch orderMsg.Task {
	case "buyOrderPlaced":
		// Process buy order
		common.BuyOrder(orderMsg.Task, UserID, orderMsg.Amount)
		log.Printf("📌 Buy Order Payment changed (User ID: %s)\n", UserID.String())
	case "sellOrderPlaced":
		common.sellOrder(orderMsg.Task, UserID, orderMsg.Amount)
		log.Printf("📌 Sell Order Payment changed (User ID: %s)\n", UserID.String())
	default:
		log.Printf("⚠️ Unknown task type: %s\n", orderMsg.Task)
		msg.Nack(false, false)
		return
	}

	msg.Ack(false)
}

// Close gracefully shuts down the consumer
func (c *RabbitMQConsumer) Close() {
	if c.Ch != nil {
		_ = c.Ch.Close()
	}
	if c.Conn != nil {
		_ = c.Conn.Close()
	}
	log.Println("🛑 RabbitMQ connection closed.")
}

func (c *RabbitMQConsumer) PublishMessage(exchange, routingKey string, message interface{}) bool {
	if c.Ch == nil {
		logger.Error("❌ RabbitMQ channel not initialized. Cannot publish.")
		return false
	}

	body, err := json.Marshal(message)
	if err != nil {
		logger.Error("❌ Failed to marshal message", "error", err)
		return false
	}

	err = c.Ch.Publish(
		exchange,
		routingKey,
		false, // Mandatory
		false, // Immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make messages persistent
		},
	)
	if err != nil {
		logger.Error("❌ Failed to publish message",
			"exchange", exchange,
			"routing_key", routingKey,
			"error", err)
		return false
	}

	logger.Info("📤 Published event",
		"exchange", exchange,
		"routing_key", routingKey,
		"message", message)
	return true
}
