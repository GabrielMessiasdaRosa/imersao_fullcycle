package transformer

import (
	"github.com/GabrielMessiasdaRosa/imersao_fullcycle/go/internal/market/dto"
	"github.com/GabrielMessiasdaRosa/imersao_fullcycle/go/internal/market/entity"
)

func TransformInput(input dto.TradeInput) *entity.Order {
	asset := entity.NewAsset(input.AssetID, input.AssetID, 1000)
	investor := entity.NewInvestor(input.InvestorID)
	order := entity.NewOrder(investor, asset, input.Shares, input.Price, input.OrderType)
	if input.CurrentShares > 0 {
		assePosition := entity.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assePosition)
	}
	return order
}
