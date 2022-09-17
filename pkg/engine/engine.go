package engine

import (
	"context"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
)

// MatchingEngine checks for matching for every added order.
// An Event will be created case necessary.
type MatchingEngine interface {
	// AddOrder adds a new order checking for matches.
	AddOrder(ctx context.Context, order entity.Order) error
	// CancelOrder remove an order by id.
	CancelOrder(ctx context.Context, orderID entity.OrderID) error
}
