package entity

type Order struct {
	ID            string         `json:"id"`
	Investor      *Investor      `json:"investor"`
	Asset         *Asset         `json:"asset"`
	Shares        int            `json:"shares"`
	PendingShares int            `json:"pending_shares"`
	Price         float64        `json:"price"`
	OrderType     string         `json:"order_type"`
	Status        string         `json:"status"`
	Transactions  []*Transaction `json:"transactions"`
}

func NewOrder(id string, investor *Investor, asset *Asset, shares int, price float64, orderType string) *Order {
	return &Order{
		ID:            id,
		Investor:      investor,
		Asset:         asset,
		Shares:        shares,
		PendingShares: shares,
		Price:         price,
		OrderType:     orderType,
		Status:        "OPEN",
		Transactions:  []*Transaction{},
	}
}
