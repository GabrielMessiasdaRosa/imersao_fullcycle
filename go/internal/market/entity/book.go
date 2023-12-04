package entity

import (
	"container/heap"
	"sync"
)

type Book struct {
	Order         []*Order        `json:"order"`
	Transactions  []*Transaction  `json:"transaction"`
	OrdersChan    chan *Order     `json:"orders_chan"` // receive orders from kafka
	OrdersChanOut chan *Order     `json:"orders_chan_out"`
	Wg            *sync.WaitGroup `json:"wg"`
}

func NewBook(orderChan chan *Order, orderChanOut chan *Order, wg *sync.WaitGroup) *Book {
	return &Book{
		Order:         []*Order{},
		Transactions:  []*Transaction{},
		OrdersChan:    orderChan,
		OrdersChanOut: orderChanOut,
		Wg:            wg,
	}
}

func (b *Book) Trade() {
	sellOrders := NewOrderQueue()
	buyOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.OrdersChan {
		switch order.OrderType {
		case "BUY": // hrdcoded mudar depois 
			b.processBuyOrder(order, sellOrders, buyOrders)
		case "SELL":
			b.processSellOrder(order, sellOrders, buyOrders)
		}
	}
}
func (book *Book) processBuyOrder(o *Order, so *OrderQueue, bo *OrderQueue) {
	bo.Push(o)
	if so.Len() > 0 && so.Orders[0].Price <= o.Price {
		s := so.Pop().(*Order)
		if s.PendingShares > 0 {
			transaction := NewTransaction(o.ID, s, o, o.Shares, s.Price)
			book.AddTransaction(transaction, book.Wg)
			s.Transacions = append(s.Transacions, transaction)
			o.Transacions = append(o.Transacions, transaction)
			book.OrdersChanOut <- s
			book.OrdersChanOut <- o
			if s.PendingShares > 0 {
				so.Push(s)
			}
		}
	}
}

func (book *Book) processSellOrder(o *Order, so *OrderQueue, bo *OrderQueue) {
	so.Push(o)
	if bo.Len() > 0 && bo.Orders[0].Price >= o.Price {
		b := bo.Pop().(*Order)
		if b.PendingShares > 0 {
			transaction := NewTransaction(o.ID, o, b, o.Shares, b.Price)
			book.AddTransaction(transaction, book.Wg)
			b.Transacions = append(b.Transacions, transaction)
			o.Transacions = append(o.Transacions, transaction)
			book.OrdersChanOut <- b
			book.OrdersChanOut <- o
			if b.PendingShares > 0 {
				bo.Push(b)
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, wg *sync.WaitGroup) {
	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares
	minShares := sellingShares

	if buyingShares < minShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellOrderPendingShares(-minShares)

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.AddBuyOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)
	transaction.CloseBuyOrder()
	transaction.CloseSellOrder()
	b.Transactions = append(b.Transactions, transaction)

	defer wg.Done()
}
