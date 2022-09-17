package engine

import (
	"context"
	"fmt"
	"sync"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
	"github.com/rodoufu/simple-orderbook/pkg/event"
)

var (
	notStartedError         = fmt.Errorf("engine not started or does not exist")
	invalidOrderAmountError = fmt.Errorf("invalid order amount")
)

type listEngine struct {
	mtx      sync.Mutex
	orders   map[entity.Side][]entity.Order
	events   chan event.Event
	orderIDs map[entity.OrderID]entity.Side
}

func (s *listEngine) Close() error {
	if s == nil {
		return nil
	}
	close(s.events)
	return nil
}

func (s *listEngine) AddOrder(ctx context.Context, order entity.Order) error {
	if s == nil {
		return notStartedError
	}
	if order.Amount == 0 {
		return invalidOrderAmountError
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()

	checkBeforeAndAfter := func(side entity.Side, before, after *entity.Order) {
		if before == nil && after == nil {
			return
		}
		if after == nil {
			s.events <- &event.TopOfBookChange{
				Side: side,
			}
		} else if before != after {
			s.events <- &event.TopOfBookChange{
				Side:          side,
				Price:         after.Price,
				TotalQuantity: after.Amount,
			}
		}
	}
	var sellTop, buyTop *entity.Order
	if len(s.orders[entity.Sell]) > 0 {
		sellTop = &s.orders[entity.Sell][len(s.orders[entity.Sell])-1]
	}
	if len(s.orders[entity.Buy]) > 0 {
		buyTop = &s.orders[entity.Buy][len(s.orders[entity.Buy])-1]
	}
	defer func() {
		var nSellTop, nBuyTop *entity.Order
		if len(s.orders[entity.Sell]) > 0 {
			nSellTop = &s.orders[entity.Sell][len(s.orders[entity.Sell])-1]
		}
		if len(s.orders[entity.Buy]) > 0 {
			nBuyTop = &s.orders[entity.Buy][len(s.orders[entity.Buy])-1]
		}

		checkBeforeAndAfter(entity.Sell, sellTop, nSellTop)
		checkBeforeAndAfter(entity.Buy, buyTop, nBuyTop)
	}()

	if _, orderExists := s.orderIDs[order.ID]; orderExists {
		return fmt.Errorf("order %v alreday exists", order.ID)
	}

	oppositeBook := s.orders[order.Side.Opposite()]
	for i := len(oppositeBook) - 1; i >= 0; i-- {
		remainingOrder, trade := order.Match(&oppositeBook[i])
		if remainingOrder == nil && trade == nil {
			break
		}
		if trade != nil {
			s.events <- &event.TradeGenerated{
				Trade: *trade,
			}
		}
		if remainingOrder == nil {
			delete(s.orderIDs, oppositeBook[i].ID)
			s.events <- &event.OrderFilled{
				Order: oppositeBook[i],
				Full:  true,
			}
			order.Amount = 0
			oppositeBook = oppositeBook[:i]
			break
		}

		if remainingOrder.Side == order.Side {
			s.events <- &event.OrderFilled{
				Order: oppositeBook[i],
				Full:  true,
			}
			oppositeBook = oppositeBook[:i]
			order = *remainingOrder
		} else {
			s.events <- &event.OrderFilled{
				Order: *remainingOrder,
				Full:  false,
			}
			oppositeBook = oppositeBook[:i]
			oppositeBook = append(oppositeBook, *remainingOrder)
			order.Amount = 0
			break
		}
	}
	s.orders[order.Side.Opposite()] = oppositeBook

	if order.Amount > 0 {
		book := s.orders[order.Side]
		book = append(book, order)
		s.orderIDs[order.ID] = order.Side
		s.events <- &event.OrderCreated{
			Order: order,
		}
		for i := len(book) - 1; i >= 1 && book[i].Less(&book[i-1]); i-- {
			book[i], book[i-1] = book[i-1], book[i]
		}

		s.orders[order.Side] = book
	}

	return nil
}

func (s *listEngine) CancelOrder(ctx context.Context, orderID entity.OrderID) error {
	if s == nil {
		return notStartedError
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()

	side, orderExists := s.orderIDs[orderID]
	if !orderExists {
		return fmt.Errorf("order %v not found", orderID)
	}

	sideOrders := s.orders[side]
	index := len(sideOrders) - 1
	for index >= 0 && sideOrders[index].ID != orderID {
		index--
	}
	if index >= 0 && sideOrders[index].ID == orderID {
		delete(s.orderIDs, orderID)
		s.events <- &event.OrderCancelled{
			Order: sideOrders[index],
		}

		copy(sideOrders[index:], sideOrders[index+1:])
		sideOrders = sideOrders[:len(sideOrders)-1]
		s.orders[side] = sideOrders

		return nil
	}

	return fmt.Errorf("order %v not found", orderID)
}

func NewListEngine() (MatchingEngine, <-chan event.Event) {
	engine := listEngine{
		mtx: sync.Mutex{},
		orders: map[entity.Side][]entity.Order{
			entity.Buy:  {},
			entity.Sell: {},
		},
		events:   make(chan event.Event, 10),
		orderIDs: map[entity.OrderID]entity.Side{},
	}
	return &engine, engine.events
}
