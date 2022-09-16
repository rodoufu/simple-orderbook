package engine

import (
	"context"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
)

type MatchingEngine interface {
	AddOrder(ctx context.Context, order entity.Order) error
	CancelOrder(ctx context.Context, orderID entity.OrderID) error
}
