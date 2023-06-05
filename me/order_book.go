package me

import (
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type OrderBook struct {
	tradePair     string
	ChNewOrder    chan QueueItem
	ChTradedOrder chan TradedOrder
	ChCancelTrade chan string
	askQueue      *OrderQueue
	bidQueue      *OrderQueue

	latestPrice  decimal.Decimal
	tradedOrders []TradedOrder // TODO: persitence storage for traded order
	sync.Mutex
}

func NewOrderBook(tradePair string) *OrderBook {
	o := &OrderBook{
		tradePair:     tradePair,
		ChNewOrder:    make(chan QueueItem),
		ChTradedOrder: make(chan TradedOrder, 10),
		ChCancelTrade: make(chan string, 10),
		askQueue:      NewQueue(),
		bidQueue:      NewQueue(),
	}

	go o.matching()
	go o.storingTradedOrder()
	return o
}

func (o *OrderBook) AskLen() int {
	return o.askQueue.Len()
}

func (o *OrderBook) BidLen() int {
	return o.bidQueue.Len()
}

func (o *OrderBook) GetTradedOrders() []TradedOrder {
	return o.tradedOrders
}

func (o *OrderBook) PushNewOrder(order QueueItem) {
	o.ChNewOrder <- order
}

func (o *OrderBook) CancelOrder(orderType OrderType, orderID string) {
	if orderType == BUYORDER {
		o.bidQueue.Remove(orderID)
	} else {
		o.askQueue.Remove(orderID)
	}
}

func (o *OrderBook) matching() {
	for {
		select {
		case newOrder := <-o.ChNewOrder:
			go o.processNewOrder(newOrder)
		default:
			o.processLimitOrder()
		}
	}
}

func (o *OrderBook) processNewOrder(order QueueItem) {
	if order.GetOrderType() == BUYORDER {
		o.bidQueue.Push(order)
	} else {
		o.askQueue.Push(order)
	}
}

func (o *OrderBook) processLimitOrder() {
	if o.askQueue == nil || o.bidQueue == nil {
		return
	} else if o.askQueue.Len() == 0 || o.bidQueue.Len() == 0 {
		return
	}

	askTopOrder := o.askQueue.Top()
	bidTopOrder := o.bidQueue.Top()
	defer func() {
		if askTopOrder.GetUnits().Equal(decimal.Zero) {
			o.askQueue.Remove(askTopOrder.GetOrderID())
		}

		if bidTopOrder.GetUnits().Equal(decimal.Zero) {
			o.bidQueue.Remove(bidTopOrder.GetOrderID())
		}
	}()

	if bidTopOrder.GetPrice().GreaterThanOrEqual(askTopOrder.GetPrice()) {
		var tradeUnits, tradePrice decimal.Decimal
		tradeUnits = minVal(askTopOrder.GetUnits(), bidTopOrder.GetUnits())
		askTopOrder.SetUnits(askTopOrder.GetUnits().Sub(tradeUnits))
		bidTopOrder.SetUnits(bidTopOrder.GetUnits().Sub(tradeUnits))

		if askTopOrder.GetCreateTime() >= bidTopOrder.GetCreateTime() {
			tradePrice = bidTopOrder.GetPrice()
		} else {
			tradePrice = askTopOrder.GetPrice()
		}
		o.sendTradeResult(askTopOrder.GetOrderID(), bidTopOrder.GetOrderID(), tradeUnits, tradePrice)
		return
	}
}

func (o *OrderBook) sendTradeResult(askOrderID, bidOrderID string, units, price decimal.Decimal) {
	o.Lock()
	defer o.Unlock()

	var tradedOrder TradedOrder
	tradedOrder.TradePair = o.tradePair
	tradedOrder.AskOrderID = askOrderID
	tradedOrder.BidOrderID = bidOrderID
	tradedOrder.createdTime = time.Now().Unix()
	tradedOrder.Units = units
	tradedOrder.Price = price
	tradedOrder.Amount = units.Mul(price)

	o.latestPrice = price
	o.ChTradedOrder <- tradedOrder
}

func (o *OrderBook) storingTradedOrder() {
	for {
		tradedOrder := <-o.ChTradedOrder
		o.tradedOrders = append(o.tradedOrders, tradedOrder)
	}
}

func (o *OrderBook) cleanAll() {
	o.askQueue.clean()
	o.bidQueue.clean()
	o.tradedOrders = []TradedOrder{}
}

func minVal(a, b decimal.Decimal) decimal.Decimal {
	if a.GreaterThan(b) {
		return b
	}
	return a
}
