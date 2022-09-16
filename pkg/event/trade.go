package event

import "github.com/rodoufu/simple-orderbook/pkg/entity"

type TradeGenerated struct {
	Event
	Trade entity.Trade
}
