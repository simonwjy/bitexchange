package me

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_OrderQueue(t *testing.T) {
	t.Run("add a new bid order", func(t *testing.T) {
		bidQueue := NewQueue()
		bidItem1 := NewBidItem("id01", decimal.NewFromFloat(1), decimal.NewFromFloat(10), decimal.Zero, 1110)
		bidItem2 := NewBidItem("id02", decimal.NewFromFloat(3), decimal.NewFromFloat(100), decimal.Zero, 1111)
		bidItem3 := NewBidItem("id03", decimal.NewFromFloat(1), decimal.NewFromFloat(200), decimal.Zero, 1112)
		bidQueue.Push(bidItem1)
		bidQueue.Push(bidItem2)
		bidQueue.Push(bidItem3)

		assert.Equal(t, bidQueue.Len(), 3, "expected depth for bid queue")
		assert.Equal(t, bidQueue.Top().GetAmount(), decimal.NewFromFloat(100), "top order amount")
		assert.Equal(t, bidQueue.Top().GetPrice(), decimal.NewFromFloat(3), "top order price")

		// remove one order
		newItem := bidQueue.Remove("id02")
		assert.Equal(t, newItem.GetOrderID(), "id02", "deleted order")
		assert.Equal(t, bidQueue.Top().GetPrice(), decimal.NewFromFloat(1), "top order price")
		assert.Equal(t, bidQueue.Top().GetAmount(), decimal.NewFromFloat(10), "top order price")
		newItem = bidQueue.Remove("id01")
		assert.Equal(t, newItem.GetOrderID(), "id01", "deleted order")
		assert.Equal(t, bidQueue.Top().GetPrice(), decimal.NewFromFloat(1), "top order price")
		assert.Equal(t, bidQueue.Top().GetAmount(), decimal.NewFromFloat(200), "top order price")
	})

	t.Run("add a new ask order", func(t *testing.T) {
		askQueue := NewQueue()
		askItem1 := NewAskItem("id01", decimal.NewFromFloat(1), decimal.NewFromFloat(10), decimal.Zero, 1110)
		askItem2 := NewAskItem("id02", decimal.NewFromFloat(3), decimal.NewFromFloat(100), decimal.Zero, 1111)
		askItem3 := NewAskItem("id03", decimal.NewFromFloat(1), decimal.NewFromFloat(200), decimal.Zero, 1112)
		askQueue.Push(askItem1)
		askQueue.Push(askItem2)
		askQueue.Push(askItem3)

		assert.Equal(t, askQueue.Len(), 3, "expected depth for ask queue")
		assert.Equal(t, askQueue.Top().GetAmount(), decimal.NewFromFloat(10), "top order amount")
		assert.Equal(t, askQueue.Top().GetPrice(), decimal.NewFromFloat(1), "top order price")

		// remove one order
		newItem := askQueue.Remove("id01")
		assert.Equal(t, newItem.GetOrderID(), "id01", "deleted order")
		assert.Equal(t, askQueue.Top().GetPrice(), decimal.NewFromFloat(1), "top order price")
		newItem = askQueue.Remove("id03")
		assert.Equal(t, newItem.GetOrderID(), "id03", "deleted order")
		assert.Equal(t, askQueue.Top().GetPrice(), decimal.NewFromFloat(3), "top order price")
	})
}
