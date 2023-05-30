package me

import (
	"sync"
)

type OrderBook struct {
	tradePair string
	askQueue  *OrderQueue
	bidQueue  *OrderQueue

	sync.Mutex
}
