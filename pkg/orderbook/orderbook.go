package orderbook

import (
	"context"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
	"github.com/rodoufu/simple-orderbook/pkg/event"
)

type OrderBook interface {
	ProcessEvent(ctx context.Context, event event.Event) error

	Bids(ctx context.Context) <-chan entity.Order
	Asks(ctx context.Context) <-chan entity.Order
	TopBid(ctx context.Context) *entity.Order
	TopAsk(ctx context.Context) *entity.Order
}
