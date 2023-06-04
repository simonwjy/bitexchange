package me

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tradeBTC = NewOrderBook("BTC_USDT")

func Test_OrderBook(t *testing.T) {
	t.Run("fully traded for the limited order", func(t *testing.T) {
		tradeBTC.cleanAll()
		bidItem1 := NewBidItem("bid1", decimal.NewFromFloat(2), decimal.Zero, decimal.NewFromFloat(10), 1114)
		bidItem2 := NewBidItem("bid2", decimal.NewFromFloat(3), decimal.Zero, decimal.NewFromFloat(10), 1115)
		tradeBTC.PushNewOrder(bidItem1)
		tradeBTC.PushNewOrder(bidItem2)

		askItem1 := NewAskItem("ask1", decimal.NewFromFloat(1), decimal.Zero, decimal.NewFromFloat(10), 1112)
		askItem2 := NewAskItem("ask2", decimal.NewFromFloat(1), decimal.Zero, decimal.NewFromFloat(10), 1113)
		tradeBTC.PushNewOrder(askItem1)
		tradeBTC.PushNewOrder(askItem2)

		time.Sleep(1 * time.Second)
		require.True(t, len(tradeBTC.tradedOrders) == 2)
		assert.Equal(t, tradeBTC.tradedOrders[0].AskOrderID, "ask1")
		assert.Equal(t, tradeBTC.tradedOrders[0].BidOrderID, "bid2")
		assert.Equal(t, tradeBTC.tradedOrders[0].Price, decimal.NewFromFloat(1))
		assert.Equal(t, tradeBTC.tradedOrders[0].Units, decimal.NewFromFloat(10))
		assert.Equal(t, tradeBTC.tradedOrders[0].Amount, decimal.NewFromFloat(10))
		assert.Equal(t, tradeBTC.tradedOrders[1].AskOrderID, "ask2")
		assert.Equal(t, tradeBTC.tradedOrders[1].BidOrderID, "bid1")
		assert.Equal(t, tradeBTC.tradedOrders[1].Price, decimal.NewFromFloat(1))
		assert.Equal(t, tradeBTC.tradedOrders[1].Units, decimal.NewFromFloat(10))
		assert.Equal(t, tradeBTC.tradedOrders[1].Amount, decimal.NewFromFloat(10))
	})
}
