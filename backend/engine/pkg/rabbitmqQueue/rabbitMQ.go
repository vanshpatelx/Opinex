// dummyengine/pkg/rabbitmqQueue/consumer.go
package rabbitmqQueue

import (
	"dummyengine/pkg/exchange"
	"encoding/json"
	"dummyengine/pkg/logger"
	"math/big"
	"time"

	"github.com/streadway/amqp"
)

// RabbitMQConsumer handles message consumption
type RabbitMQQueue struct {
	Conn     *amqp.Connection
	Ch       *amqp.Channel
	Queue    string
	URL      string
	Exchange *exchange.Exchange
}

// EventMessage represents the structure received from RabbitMQ
type EventMessage struct {
	Task          string `json:"task"`               // "Order" || "CreateEvent" || "Settlement"
	ID            string `json:"eventId"`            // EventID
	OrderID       string `json:"orderId,omitempty"`  // OrderID
	OrderPrice    int    `json:"price,omitempty"`    // OrderPrice
	OrderUserID   string `json:"userId,omitempty"`   // OrderUserID
	OrderQuantity int    `json:"quantity,omitempty"` // OrderQuantity
	Type          string `json:"type,omitempty"`     // "BUY" || "SELL"
}

// NewRabbitMQConsumer initializes a new consumer
func NewRabbitMQQueue(url, queue, tradeQueue, priceQueue, orderBookQueue string) *RabbitMQQueue {
	return &RabbitMQQueue{
		Queue:    queue,
		URL:      url,
		Exchange: nil, // Initialize to nil, will be set after connection
	}
}

// Connect establishes the RabbitMQ connection
func (c *RabbitMQQueue) Connect() {
	var err error

	// Attempt connection
	c.Conn, err = amqp.Dial(c.URL)
	if err != nil {
		logger.Fatal("❌ Failed to connect to RabbitMQ", "error", err)
	}

	c.Ch, err = c.Conn.Channel()
	if err != nil {
		logger.Fatal("❌ Failed to open a channel", "error", err)
	}

	logger.Info("✅ Connected to RabbitMQ", "url", c.URL, "queue", c.Queue)

	// Auto-reconnect on failure
	c.Conn.NotifyClose(make(chan *amqp.Error))
	c.reconnect()

	// Start consuming
	c.Consume()
}

// Reconnect handles RabbitMQ reconnection logic
func (c *RabbitMQQueue) reconnect() {
	for {
		logger.Warn("🔄 Reconnecting to RabbitMQ...")
		time.Sleep(5 * time.Second)

		var err error
		c.Conn, err = amqp.Dial(c.URL)
		if err == nil {
			c.Ch, err = c.Conn.Channel()
			if err == nil {
				logger.Info("✅ Reconnected to RabbitMQ")
				c.Exchange = exchange.NewExchange(c.Ch, "trade_queue", "price_queue", "orderBook_queue")
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
		logger.Fatal("❌ Failed to register a consumer", "error", err)
	}

	logger.Info("📥 Consuming messages", "queue", c.Queue)

	for msg := range msgs {
		c.processMessage(msg)
	}
}

func (c *RabbitMQQueue) processMessage(msg amqp.Delivery) {
	var orderMsg EventMessage
	if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
		logger.Error("❌ Failed to parse message", "error", err)
		msg.Nack(false, false) // Reject message
		return
	}

	// Convert string values to big.Int
	ID := new(big.Int)
	if err := ID.UnmarshalText([]byte(orderMsg.ID)); err != nil {
		logger.Error("❌ Failed to convert orderMsg.ID to big.Int", "error", err, "event", orderMsg)
		return
	}

	switch orderMsg.Task {
	case "CreateEvent":
		logger.Info("📌 Creating event", "event_id", ID.String())
		c.Exchange.AddEvent(ID)

	case "Settlement":
		logger.Info("💰 Processing settlement", "event_id", ID.String())
		c.Exchange.Settlement(ID)

	case "Order":
		orderID := new(big.Int)
		if err := orderID.UnmarshalText([]byte(orderMsg.OrderID)); err != nil {
			logger.Error("❌ Failed to convert orderMsg.OrderID to big.Int", "error", err, "event", orderMsg)
			return
		}

		OrderUserID := new(big.Int)
		if err := OrderUserID.UnmarshalText([]byte(orderMsg.OrderUserID)); err != nil {
			logger.Error("❌ Failed to convert orderMsg.OrderUserID to big.Int", "error", err, "event", orderMsg)
			return
		}

		switch orderMsg.Type {
		case "BUY":
			c.Exchange.AddBuyOrder(ID, orderID, orderMsg.OrderPrice, orderMsg.OrderQuantity, OrderUserID)
		case "SELL":
			c.Exchange.AddSellOrder(ID, orderID, orderMsg.OrderPrice, orderMsg.OrderQuantity, OrderUserID)
		default:
			logger.Warn("⚠️ Unknown order type", "order_type", orderMsg.Type)
			msg.Nack(false, false)
			return
		}

	default:
		logger.Warn("⚠️ Unknown task type", "task_type", orderMsg.Task)
		msg.Nack(false, false)
		return
	}

	msg.Ack(false) // Acknowledge message
}
