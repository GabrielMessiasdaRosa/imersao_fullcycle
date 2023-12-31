package entity

import (
	"container/heap"
	"fmt"
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
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)
	/* 	sellOrders := NewOrderQueue()
	   	buyOrders := NewOrderQueue() */

	/* heap.Init(buyOrders)
	heap.Init(sellOrders) */

	for order := range b.OrdersChan {
		fmt.Println("order received")
		asset := order.Asset.ID

		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}

		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}
		switch order.OrderType {
		case "BUY": // hrdcoded mudar depois
			b.executeBuyOrderMatchingSellOrders(order, sellOrders[asset], buyOrders[asset])
		case "SELL":
			b.executeSellOrderMatchingBuyOrders(order, sellOrders[asset], buyOrders[asset])
		}
	}
}

func (book *Book) executeBuyOrderMatchingSellOrders(buyOrder *Order, sellOrders *OrderQueue, buyOrders *OrderQueue) {
	buyOrders.Push(buyOrder)
	sellOrdersAvailable := sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= buyOrder.Price
	if sellOrdersAvailable {
		matchedSellOrder := sellOrders.Pop().(*Order)
		hasPendingShares := matchedSellOrder.PendingShares > 0
		if hasPendingShares {
			transaction := NewTransaction(buyOrder.ID, matchedSellOrder, buyOrder, buyOrder.Shares, matchedSellOrder.Price)
			book.AddTransaction(transaction, book.Wg)
			matchedSellOrder.Transactions = append(matchedSellOrder.Transactions, transaction)
			buyOrder.Transactions = append(buyOrder.Transactions, transaction)
			book.OrdersChanOut <- matchedSellOrder
			book.OrdersChanOut <- buyOrder
			if hasPendingShares {
				sellOrders.Push(matchedSellOrder)
			}
		}
	}
}

func (book *Book) executeSellOrderMatchingBuyOrders(sellOrder *Order, sellOrders *OrderQueue, buyOrders *OrderQueue) {
	sellOrders.Push(sellOrder)
	buyOrdersAvailable := buyOrders.Len() > 0 && buyOrders.Orders[0].Price >= sellOrder.Price
	fmt.Println("buyOrdersAvailable", buyOrdersAvailable)
	if buyOrdersAvailable {
		matchedBuyOrder := buyOrders.Pop().(*Order)
		fmt.Println("matchedBuyOrder", matchedBuyOrder.PendingShares > 0)
		hasPendingShares := matchedBuyOrder.PendingShares > 0
		if hasPendingShares {
			transaction := NewTransaction(sellOrder.ID, sellOrder, matchedBuyOrder, sellOrder.Shares, matchedBuyOrder.Price)
			fmt.Println("Transaction craeted")
			book.AddTransaction(transaction, book.Wg)
			matchedBuyOrder.Transactions = append(matchedBuyOrder.Transactions, transaction)
			sellOrder.Transactions = append(sellOrder.Transactions, transaction)
			book.OrdersChanOut <- matchedBuyOrder
			book.OrdersChanOut <- sellOrder
			if hasPendingShares {
				buyOrders.Push(matchedBuyOrder)
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
	fmt.Println("Transaction added")
	defer wg.Done()
}
