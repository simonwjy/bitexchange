package me

import (
	"container/heap"
	"sync"

	"github.com/shopspring/decimal"
)

type QueueItem interface {
	GetOrderID() string
	SetIndex(index int)
	GetIndex() int
	GetPrice() decimal.Decimal
	GetAmount() decimal.Decimal
	GetUnits() decimal.Decimal
	GetCreateTime() int64
	Less(item QueueItem) bool
}

type PriorityQueue []QueueItem

func (p PriorityQueue) Len() int { return len(p) }

func (p PriorityQueue) Less(i, j int) bool {
	return p[i].Less(p[j])
}

func (p PriorityQueue) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].SetIndex(i)
	p[j].SetIndex(j)
}

func (p *PriorityQueue) Pop() interface{} {
	n := len(*p)
	v := (*p)[n-1]
	v.SetIndex(-1)
	*p = (*p)[:n-1]
	return v
}

func (p *PriorityQueue) Push(v interface{}) {
	n := len(*p)
	v.(QueueItem).SetIndex(n)
	*p = append(*p, v.(QueueItem))
}

func NewQueue() *OrderQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	return &OrderQueue{
		pq: &pq,
		m:  make(map[string]*QueueItem),
	}
}

type OrderQueue struct {
	pq *PriorityQueue
	m  map[string]*QueueItem

	sync.Mutex
}

func (o *OrderQueue) Len() int {
	return o.pq.Len()
}

func (o *OrderQueue) Push(item QueueItem) bool {
	o.Lock()
	defer o.Unlock()

	orderID := item.GetOrderID()
	if _, ok := o.m[orderID]; ok {
		return true
	}

	heap.Push(o.pq, item)
	o.m[orderID] = &item
	return false
}

func (o *OrderQueue) Get(index int) QueueItem {
	n := o.pq.Len()
	if n <= index {
		return nil
	}
	return (*o.pq)[index]
}

func (o *OrderQueue) Top() QueueItem {
	return (*o.pq)[0]
}

func (o *OrderQueue) Remove(orderID string) QueueItem {
	o.Lock()
	defer o.Unlock()

	oldItem, ok := o.m[orderID]
	if !ok {
		return nil
	}

	item := heap.Remove(o.pq, (*oldItem).GetIndex())
	delete(o.m, orderID)
	return item.(QueueItem)
}
