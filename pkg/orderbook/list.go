package orderbook

import (
	"context"
	"fmt"
	"sync"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
	"github.com/rodoufu/simple-orderbook/pkg/event"
)

var (
	notStartedError = fmt.Errorf("orderbook not started or does not exist")
)

type listOrderBook struct {
	mtx    map[entity.Side]*sync.RWMutex
	orders map[entity.Side][]entity.Order
}

func (l *listOrderBook) TopBid(ctx context.Context) *BookLevel {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	levels := l.Bids(ctx)
	done := ctx.Done()
	select {
	case <-done:
		return nil
	case level, ok := <-levels:
		if !ok {
			return nil
		}
		return &level
	}
}

func (l *listOrderBook) TopAsk(ctx context.Context) *BookLevel {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	levels := l.Asks(ctx)
	done := ctx.Done()
	select {
	case <-done:
		return nil
	case level, ok := <-levels:
		if !ok {
			return nil
		}
		return &level
	}
}

func (l *listOrderBook) ProcessEvent(ctx context.Context, evt event.Event) error {
	if l == nil {
		return notStartedError
	}

	switch it := evt.(type) {
	case *event.TradeGenerated:
	case *event.OrderCancelled:
		return l.cancelOrder(ctx, it.Order.ID, it.Order.Side)
	case *event.OrderCreated:
		return l.addOrder(ctx, it.Order)
	case *event.OrderUpdated:
		return l.updateOrder(ctx, it.Order)
	case *event.OrderFilled:
		if it.Full {
			return l.cancelOrder(ctx, it.Order.ID, it.Order.Side)
		} else {
			return l.updateOrder(ctx, it.Order)
		}
	default:
		return fmt.Errorf("unexpected event: %v", evt)
	}
	return nil
}

func (l *listOrderBook) getLevel(ctx context.Context, side entity.Side) <-chan BookLevel {
	l.mtx[side].RLock()
	resp := make(chan BookLevel)
	if l == nil {
		close(resp)
		return resp
	}

	go func() {
		defer l.mtx[side].RUnlock()
		defer close(resp)
		done := ctx.Done()

		for i := len(l.orders[side]) - 1; i >= 0; i-- {
			order := l.orders[side][i]
			level := BookLevel{
				Side:  side,
				Price: order.Price,
			}

			for ; i >= 0 && order.Price == l.orders[side][i].Price; i-- {
				level.TotalQuantity += l.orders[side][i].Amount
			}
			if i >= 0 && order.Price != l.orders[side][i].Price {
				i++
			}

			select {
			case <-done:
				return
			case resp <- level:
			}
		}
	}()

	return resp

}

func (l *listOrderBook) Bids(ctx context.Context) <-chan BookLevel {
	return l.getLevel(ctx, entity.Buy)

}

func (l *listOrderBook) Asks(ctx context.Context) <-chan BookLevel {
	return l.getLevel(ctx, entity.Sell)
}

func (l *listOrderBook) cancelOrder(ctx context.Context, orderID entity.OrderID, side entity.Side) error {
	if l == nil {
		return notStartedError
	}
	l.mtx[side].Lock()
	defer l.mtx[side].Unlock()

	sideOrders := l.orders[side]
	index := len(sideOrders) - 1
	for index >= 0 && sideOrders[index].ID != orderID {
		index--
	}
	if index >= 0 && sideOrders[index].ID == orderID {
		copy(sideOrders[index:], sideOrders[index+1:])
		sideOrders = sideOrders[:len(sideOrders)-1]
		l.orders[side] = sideOrders

		return nil
	}

	return fmt.Errorf("order %v not found", orderID)
}

func (l *listOrderBook) addOrder(ctx context.Context, order entity.Order) error {
	if l == nil {
		return notStartedError
	}
	l.mtx[order.Side].Lock()
	defer l.mtx[order.Side].Unlock()

	sideOrders := l.orders[order.Side]
	sideOrders = append(sideOrders, order)

	for i := len(sideOrders) - 1; i >= 1 && sideOrders[i].Less(&sideOrders[i-1]); i-- {
		sideOrders[i], sideOrders[i-1] = sideOrders[i-1], sideOrders[i]
	}

	l.orders[order.Side] = sideOrders

	return nil
}

func (l *listOrderBook) updateOrder(ctx context.Context, order entity.Order) error {
	if l == nil {
		return notStartedError
	}
	l.mtx[order.Side].Lock()
	defer l.mtx[order.Side].Unlock()

	sideOrders := l.orders[order.Side]
	found := false
	for i, it := range sideOrders {
		if order.ID == it.ID {
			sideOrders[i] = order
			found = true
		}
	}
	l.orders[order.Side] = sideOrders

	if !found {
		return fmt.Errorf("order not found: %v", order.ID)
	}
	return nil
}

func NewListOrderBook() OrderBook {
	return &listOrderBook{
		mtx: map[entity.Side]*sync.RWMutex{
			entity.Buy:  {},
			entity.Sell: {},
		},
		orders: map[entity.Side][]entity.Order{
			entity.Buy:  {},
			entity.Sell: {},
		},
	}
}
