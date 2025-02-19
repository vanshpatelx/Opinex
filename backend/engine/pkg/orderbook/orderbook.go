// dummyengine/pkg/orderbook/orderbook.go
package orderbook

import (
	"container/heap"
	"dummyengine/pkg/customheap"
	"dummyengine/pkg/pricelevel"
	"dummyengine/pkg/uniqueid"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/big"
	"time"
)

type OrderBook struct {
	BuyOrders      *customheap.BuyOrderBook
	SellOrders     *customheap.SellOrderBook
	Ch             *amqp.Channel
	TradeQueue     string
	PriceQueue     string
	OrderBookQueue string
}

type TradeMessage struct {
	ID        *big.Int `json:"id"`
	OrderID   *big.Int `json:"order_id"`
	Price     int      `json:"price"`
	Quantity  int      `json:"quantity"`
	Timestamp int64    `json:"timestamp"`
}

type PriceUpdate struct {
	Price int `json:"price"`
}

func NewOrderBook(ch *amqp.Channel, tradeQueue, priceQueue, orderBookQueue string) *OrderBook {
	buyHeap := &customheap.BuyOrderBook{}
	sellHeap := &customheap.SellOrderBook{}
	heap.Init(buyHeap)
	heap.Init(sellHeap)

	return &OrderBook{
		BuyOrders:      buyHeap,
		SellOrders:     sellHeap,
		Ch:             ch,
		TradeQueue:     tradeQueue,
		PriceQueue:     priceQueue,
		OrderBookQueue: orderBookQueue,
	}
}

func (ob *OrderBook) AddBuyOrder(order *pricelevel.Order) {
	for _, level := range ob.BuyOrders.CommonHeap {
		if level.Price == order.Price {
			level.Orders = append(level.Orders, order)
			level.Quantity += order.Quantity
			heap.Fix(ob.BuyOrders, level.Index)
			ob.MatchOrders()
			ob.publishPriceUpdate(level.Price)
			ob.publishOrderBook()
			return
		}
	}

	newLevel := &pricelevel.PriceLevel{
		Price:    order.Price,
		Quantity: order.Quantity,
		Orders:   []*pricelevel.Order{order},
	}

	heap.Push(ob.BuyOrders, newLevel)
	ob.MatchOrders()
	ob.publishPriceUpdate(newLevel.Price)
	ob.publishOrderBook()
}

func (ob *OrderBook) AddSellOrder(order *pricelevel.Order) {
	for _, level := range ob.SellOrders.CommonHeap {
		if level.Price == order.Price {
			level.Orders = append(level.Orders, order)
			level.Quantity += order.Quantity
			heap.Fix(ob.SellOrders, level.Index)
			ob.MatchOrders()
			ob.publishPriceUpdate(level.Price)
			ob.publishOrderBook()
			return
		}
	}

	newLevel := &pricelevel.PriceLevel{
		Price:    order.Price,
		Quantity: order.Quantity,
		Orders:   []*pricelevel.Order{order},
	}

	heap.Push(ob.SellOrders, newLevel)
	ob.MatchOrders()
	ob.publishPriceUpdate(newLevel.Price)
	ob.publishOrderBook()
}

func (ob *OrderBook) GetTopBuyOrder() *pricelevel.PriceLevel {
	if len(ob.BuyOrders.CommonHeap) > 0 {
		return ob.BuyOrders.CommonHeap[0]
	}
	return nil
}

func (ob *OrderBook) GetTopSellOrder() *pricelevel.PriceLevel {
	if len(ob.SellOrders.CommonHeap) > 0 {
		return ob.SellOrders.CommonHeap[0]
	}
	return nil
}

func (ob *OrderBook) GetAllBuyOrders() {
	var totalOrders int = 0
	if len(ob.BuyOrders.CommonHeap) > 0 {
		for i := range ob.BuyOrders.CommonHeap {
			orders := ob.BuyOrders.CommonHeap[i].Orders
			fmt.Printf("Buy Order: Total Q: %v  Price: %v \n", ob.BuyOrders.CommonHeap[i].Quantity, ob.BuyOrders.CommonHeap[i].Price)
			fmt.Printf("Numbers of Unique Orders: %v\n", len(orders))

			totalOrders += len(orders)
			// for _, order := range orders {
			// 	fmt.Printf("%v %v\n", order.Price, order.Quantity)
			// }
		}
	}
	fmt.Printf("Total Orders: %v\n", totalOrders)
}

func (ob *OrderBook) GetAllSellOrders() {
	var totalOrders int = 0
	if len(ob.SellOrders.CommonHeap) > 0 {
		for i := range ob.SellOrders.CommonHeap {
			orders := ob.SellOrders.CommonHeap[i].Orders
			fmt.Printf("Sell Order: Total Q: %v  Price: %v \n", ob.SellOrders.CommonHeap[i].Quantity, ob.SellOrders.CommonHeap[i].Price)
			fmt.Printf("Numbers of Unique Orders: %v\n", len(orders))

			totalOrders += len(orders)
			// for id, order := range orders {
			// 	fmt.Printf("%v %v %v\n", id, order.Price, order.Quantity)
			// }
		}
	}
	fmt.Printf("Total Orders: %v\n", totalOrders)
}

func (ob *OrderBook) MatchOrders() {
	for {
		topBuy := ob.GetTopBuyOrder()
		topSell := ob.GetTopSellOrder()

		if topBuy == nil || topSell == nil || topBuy.Price < topSell.Price {
			break
		}

		for len(topBuy.Orders) > 0 && len(topSell.Orders) > 0 {
			buyOrder := topBuy.Orders[0]
			sellOrder := topSell.Orders[0]

			matchQty := min(buyOrder.Quantity, sellOrder.Quantity)
			fmt.Printf("Matched Order: Price %v, Quantity %v Buyer: %v Seller: %v\n", topSell.Price, matchQty, buyOrder.ID, sellOrder.ID)
			buyerSideTrade := TradeMessage{
				ID:        uniqueid.GenerateBaseId(),
				OrderID:   buyOrder.ID,
				Price:     topSell.Price,
				Quantity:  matchQty,
				Timestamp: time.Now().Unix(),
			}
			ob.publishTrade(buyerSideTrade)

			sellerSideTrade := TradeMessage{
				ID:        uniqueid.GenerateBaseId(),
				OrderID:   sellOrder.ID,
				Price:     topSell.Price,
				Quantity:  matchQty,
				Timestamp: time.Now().Unix(),
			}

			ob.publishTrade(sellerSideTrade)

			buyOrder.Quantity -= matchQty
			sellOrder.Quantity -= matchQty
			topBuy.Quantity -= matchQty
			topSell.Quantity -= matchQty

			if buyOrder.Quantity == 0 {
				topBuy.Orders = topBuy.Orders[1:]
			}
			if sellOrder.Quantity == 0 {
				topSell.Orders = topSell.Orders[1:]
			}
		}

		if len(topBuy.Orders) == 0 {
			heap.Pop(ob.BuyOrders)
		} else {
			heap.Fix(ob.BuyOrders, topBuy.Index)
		}

		if len(topSell.Orders) == 0 {
			heap.Pop(ob.SellOrders)
		} else {
			heap.Fix(ob.SellOrders, topSell.Index)
		}
	}
}

func (ob *OrderBook) PublishMessage(exchange, routingKey string, message interface{}) {
	if ob.Ch == nil {
		log.Printf("❌ RabbitMQ channel not initialized. Cannot publish.")
		return
	}

	body, err := json.Marshal(message)
	if err != nil {
		log.Printf("❌ Failed to marshal message: %v", err)
		return
	}

	err = ob.Ch.Publish(
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
		log.Printf("❌ Failed to publish message to %s with key '%s': %v", exchange, routingKey, err)
		return
	}

	log.Printf("📤 Published event to %s with key '%s': %v", exchange, routingKey, message)
}

func (ob *OrderBook) publishOrderBook() {
	// Extract Price and Quantity for Buy Levels
	var buyLevels []struct {
		Price    int `json:"Price"`
		Quantity int `json:"Quantity"`
	}
	for _, level := range ob.BuyOrders.CommonHeap {
		buyLevels = append(buyLevels, struct {
			Price    int `json:"Price"`
			Quantity int `json:"Quantity"`
		}{
			Price:    level.Price,
			Quantity: level.Quantity,
		})
	}

	// Extract Price and Quantity for Sell Levels
	var sellLevels []struct {
		Price    int `json:"Price"`
		Quantity int `json:"Quantity"`
	}
	for _, level := range ob.SellOrders.CommonHeap {
		sellLevels = append(sellLevels, struct {
			Price    int `json:"Price"`
			Quantity int `json:"Quantity"`
		}{
			Price:    level.Price,
			Quantity: level.Quantity,
		})
	}

	// Construct the message
	orderBookMsg := struct {
		BuyLevels []struct {
			Price    int `json:"Price"`
			Quantity int `json:"Quantity"`
		} `json:"buy_levels"`
		SellLevels []struct {
			Price    int `json:"Price"`
			Quantity int `json:"Quantity"`
		} `json:"sell_levels"`
	}{
		BuyLevels:  buyLevels,
		SellLevels: sellLevels,
	}

	// Publish the message
	ob.PublishMessage("order_book_exchange", "order_book.update", orderBookMsg)
}

// PublishPriceUpdate publishes price changes
func (ob *OrderBook) publishPriceUpdate(price int) {
	ob.PublishMessage("price_exchange", "price.update", PriceUpdate{Price: price})
}

// PublishTrade publishes trade events
func (ob *OrderBook) publishTrade(trade TradeMessage) {
	ob.PublishMessage("trade_exchange", "trade.executed", trade)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
