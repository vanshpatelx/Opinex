package exchange

import (
	"dummyengine/pkg/orderbook"
	"math/big"
	"github.com/streadway/amqp"
	"dummyengine/pkg/logger"
)

type Exchange struct {
	OrderBooks     map[string]*orderbook.OrderBook
	TradeQueue     string
	PriceQueue     string
	OrderBookQueue string
	Ch             *amqp.Channel
}

func NewExchange(ch *amqp.Channel, tradeQueue, priceQueue, orderBookQueue string) *Exchange {
	return &Exchange{
		OrderBooks:     make(map[string]*orderbook.OrderBook),
		TradeQueue:     tradeQueue,
		PriceQueue:     priceQueue,
		OrderBookQueue: orderBookQueue,
		Ch:             ch,
	}
}

func (e *Exchange)AddEvent(eventID *big.Int) {
	eventKey := eventID.String()
	if _, exists := e.OrderBooks[eventKey]; exists{
		logger.Warn("Event already exists", "event_key", eventKey)
		return
	}

	// create orderbook
	e.OrderBooks[eventKey] = orderbook.NewOrderBook(e.Ch, eventID, e.TradeQueue, e.PriceQueue, e.OrderBookQueue)
	logger.Info("New OrderBook created", "event_key", eventKey)
}

func (e *Exchange)AddBuyOrder(eventID *big.Int, orderID *big.Int, orderPrice int, orderQuantity int, orderUserID *big.Int){
	eventKey := eventID.String()
	orderBook, exists := e.OrderBooks[eventKey]
	if !exists {
		logger.Error("OrderBook not found", "event_key", eventKey)
		return
	}

	orderBook.AddBuyOrder(orderID, orderPrice, orderQuantity, orderUserID)
	logger.Info("📦 Added buy order", "order_id", orderID.String(), "event_key", eventKey)

}
func (e *Exchange)AddSellOrder(eventID *big.Int, orderID *big.Int, orderPrice int, orderQuantity int, orderUserID *big.Int){
	eventKey := eventID.String()
	orderBook, exists := e.OrderBooks[eventKey]
	if !exists {
		logger.Warn("OrderBook not found", "event_key", eventKey)
		return
	}

	orderBook.AddSellOrder(orderID, orderPrice, orderQuantity, orderUserID)
	logger.Info("📦 Added sell order", "order_id", orderID.String(), "event_key", eventKey)
}

func (e *Exchange)Settlement(eventID *big.Int){
	
}
