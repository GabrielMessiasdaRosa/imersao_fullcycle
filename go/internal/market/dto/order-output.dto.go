package dto

type OrderOutput struct {
	OrderID           string               `json:"order_id"`
	InvestorID        string               `json:"investor_id"`
	AssetID           string               `json:"asset_id"`
	OrderType         string               `json:"order_type"`
	Status            string               `json:"status"`
	Partial           int                  `json:"partial"`
	Shares            int                  `json:"shares"`
	TransactionOutput []*TransactionOutput `json:"transaction_output"`
}
