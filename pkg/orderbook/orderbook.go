package orderbook

import (
	"context"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
	"github.com/rodoufu/simple-orderbook/pkg/event"
)

// OrderBook models another service that would implement this functions.
type OrderBook interface {
	// ProcessEvent process the events produced by the MatchingEngine so this part can be consistent.
	ProcessEvent(ctx context.Context, event event.Event) error

	// Bids returns the buy orders for the book.
	Bids(ctx context.Context) <-chan entity.Order
	// Asks returns the sell orders for the book.
	Asks(ctx context.Context) <-chan entity.Order
	// TopBid gives the top buy order.
	TopBid(ctx context.Context) *entity.Order
	// TopAsk gives the top sell order.
	TopAsk(ctx context.Context) *entity.Order
}
