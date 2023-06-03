package me

import "github.com/shopspring/decimal"

type TradedOrder struct {
	TradePair   string
	AskOrderID  string
	BidOrderID  string
	Amount      decimal.Decimal
	Units       decimal.Decimal
	Price       decimal.Decimal
	createdTime int64
}
