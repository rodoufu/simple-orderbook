package event

import (
	"fmt"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
)

// TradeGenerated is emitted when a match is found.
type TradeGenerated struct {
	Event
	Trade entity.Trade
}

func (tg *TradeGenerated) Output() string {
	if tg == nil {
		return ""
	}
	return fmt.Sprintf(
		"T, %v, %v, %v, %v, %v, %v",
		tg.Trade.BuyUserID, tg.Trade.BuyOrderID, tg.Trade.SellUserID, tg.Trade.SellOrderID,
		tg.Trade.Price, tg.Trade.Amount,
	)
}
