package orderbook

import "github.com/rodoufu/simple-orderbook/pkg/entity"

type BookLevel struct {
	Side          entity.Side
	Price         uint64
	TotalQuantity uint64
}
