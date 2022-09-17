package event

import "github.com/rodoufu/simple-orderbook/pkg/entity"

// TradeGenerated is emitted when a match is found.
type TradeGenerated struct {
	Event
	Trade entity.Trade
}
