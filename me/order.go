package me

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderType int

const (
	BUYORDER  OrderType = 0
	SELLORDER OrderType = 0
)

type Order struct {
	orderID     string
	index       int
	price       decimal.Decimal
	amount      decimal.Decimal
	units       decimal.Decimal
	createdTime int64
}

func (o *Order) GetOrderID() string {
	return o.orderID
}

func (o *Order) SetIndex(index int) {
	o.index = index
}

func (o *Order) GetIndex() int {
	return o.index
}

func (o *Order) GetPrice() decimal.Decimal {
	return o.price
}

func (o *Order) GetAmount() decimal.Decimal {
	return o.amount
}

func (o *Order) GetUnits() decimal.Decimal {
	return o.units
}

func (o *Order) GetCreateTime() int64 {
	return o.createdTime
}

type AskItem struct {
	Order
}

func (a *AskItem) Less(item QueueItem) bool {
	return a.price.Cmp(item.GetPrice()) == -1 || (a.price.Cmp(item.GetPrice()) == 0 && a.createdTime < item.GetCreateTime())
}

func (a *AskItem) GetOrderType() OrderType {
	return SELLORDER
}

func NewAskItem(orderID string, price, amount, units decimal.Decimal) *AskItem {
	return &AskItem{
		Order: Order{
			orderID:     orderID,
			price:       price,
			amount:      amount,
			units:       units,
			createdTime: time.Now().UnixMicro(),
		},
	}
}

type BidItem struct {
	Order
}

func (b *BidItem) Less(item QueueItem) bool {
	return b.price.Cmp(item.GetPrice()) == 1 || (b.price.Cmp(item.GetPrice()) == 0 && b.createdTime < item.GetCreateTime())
}

func (b *BidItem) GetOrderType() OrderType {
	return BUYORDER
}

func NewBidItem(orderID string, price, amount, units decimal.Decimal) *BidItem {
	return &BidItem{
		Order: Order{
			orderID:     orderID,
			price:       price,
			amount:      amount,
			units:       units,
			createdTime: time.Now().UnixNano(),
		},
	}
}
