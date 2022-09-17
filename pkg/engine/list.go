package engine

import (
	"context"
	"fmt"
	"sync"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
	"github.com/rodoufu/simple-orderbook/pkg/event"
	"github.com/rodoufu/simple-orderbook/pkg/io"
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

func (s *listEngine) ProcessTransaction(ctx context.Context, transaction io.Transaction) error {
	switch t := transaction.(type) {
	case io.NewOrderTransaction:
		return s.AddOrder(ctx, t.Order)
	case io.CancelOrderTransaction:
		return s.CancelOrder(ctx, t.OrderID)
	case io.ErrorTransaction:
		return t.Err
	case io.FlushAllOrdersTransaction:
		s.mtx.Lock()
		s.orders = map[entity.Side][]entity.Order{
			entity.Buy:  {},
			entity.Sell: {},
		}
		s.orderIDs = map[entity.OrderID]entity.Side{}
		s.mtx.Unlock()
		return nil
	default:
		return fmt.Errorf("problem identifying transaction: %v", transaction)
	}
}

func (s *listEngine) Close() error {
	if s == nil {
		return nil
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	close(s.events)
	return nil
}

func (s *listEngine) checkBeforeAndAfter(side entity.Side, before, after *entity.Order) {
	if before == nil && after == nil {
		return
	}
	if before == nil {
		s.events <- &event.TopOfBookChange{
			Side:          side,
			Price:         after.Price,
			TotalQuantity: after.Amount,
		}
	} else if after == nil {
		s.events <- &event.TopOfBookChange{
			Side: side,
		}
	} else if before.Price != after.Price || before.Amount != after.Amount || before.ID != after.ID {
		top := event.TopOfBookChange{
			Side:          side,
			Price:         after.Price,
			TotalQuantity: after.Amount,
		}

		for i := len(s.orders[side]) - 2; i >= 0 && s.orders[side][i].Price == top.Price; i-- {
			top.TotalQuantity += s.orders[side][i].Amount
		}

		s.events <- &top
	}
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

	before := map[entity.Side]*entity.Order{}
	for _, side := range []entity.Side{entity.Buy, entity.Sell} {
		if len(s.orders[side]) > 0 {
			before[side] = &s.orders[side][len(s.orders[side])-1]
		}
	}
	defer func() {
		for _, side := range []entity.Side{entity.Buy, entity.Sell} {
			if len(s.orders[side]) > 0 {
				s.checkBeforeAndAfter(side, before[side], &s.orders[side][len(s.orders[side])-1])
			} else {
				s.checkBeforeAndAfter(side, before[side], nil)
			}
		}
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

	before := map[entity.Side]*entity.Order{}
	if len(s.orders[side]) > 0 {
		before[side] = &s.orders[side][len(s.orders[side])-1]
	}
	defer func() {
		if len(s.orders[side]) > 0 {
			s.checkBeforeAndAfter(side, before[side], &s.orders[side][len(s.orders[side])-1])
		} else {
			s.checkBeforeAndAfter(side, before[side], nil)
		}
	}()

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
