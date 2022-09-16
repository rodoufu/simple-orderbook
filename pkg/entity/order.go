package entity

import (
	"fmt"
	"time"
)

// Side represents Sell and Buy side of the book.
type Side uint8

const (
	// InvalidSide represents the invalid initial state.
	InvalidSide Side = iota
	Buy         Side = iota
	Sell        Side = iota
)

// Opposite is used to get the other side, so the matching engine code can be generic.
func (s Side) Opposite() Side {
	switch s {
	case Buy:
		return Sell
	case Sell:
		return Buy
	case InvalidSide:
		panic(fmt.Sprintf("invalid side: %v", s))
	default:
		panic(fmt.Sprintf("invalid side: %v", s))
	}
}

func (s Side) String() string {
	switch s {
	case Buy:
		return "buy"
	case Sell:
		return "sell"
	case InvalidSide:
		return "invalid side"
	default:
		return fmt.Sprintf("invalid side (%v)", uint8(s))
	}
}

// OrderID represents the type used of orders identification.
type OrderID uint64

// UserID represents the type used of users identification.
type UserID uint64

// Order represents each order placed.
type Order struct {
	// Amount is how much the client wants to buy.
	Amount uint64
	// Price is how much the client is willing to pay.
	Price uint64
	// ID is the identification of the order.
	ID OrderID
	// Side of the book it will sit.
	Side Side
	// User identifies the user that placed the order.
	User UserID
	// Timestamp for when the order was generated.
	Timestamp time.Time
}

// Less checks if the current order should appear before the other one in the book.
func (o *Order) Less(other *Order) bool {
	if o.Side != other.Side {
		panic(fmt.Sprintf("cannot compare orders of different sides"))
	}
	switch o.Side {
	case Buy:
		return o.Price < other.Price || (o.Price == other.Price && o.Timestamp.Before(other.Timestamp))
	case Sell:
		return o.Price > other.Price || (o.Price == other.Price && o.Timestamp.Before(other.Timestamp))
	case InvalidSide:
		panic(fmt.Sprintf("invalid side: %v", o.Side))
	default:
		panic(fmt.Sprintf("invalid side: %v", o.Side))
	}
}

// Match process the matching between two orders.
// It returns the remaining order in case there is something left and the generated trade.
func (o *Order) Match(other *Order) (*Order, *Trade) {
	if o == nil || other == nil || o.Side.Opposite() != other.Side {
		return nil, nil
	}
	switch o.Side {
	case Buy:
		if o.Price >= other.Price {
			if o.Amount == other.Amount {
				return nil, &Trade{
					TakeOrderID:  o.ID,
					MakerOrderID: other.ID,
					Amount:       o.Amount,
					Price:        other.Price,
					Timestamp:    time.Now(),
				}
			} else if o.Amount > other.Amount {
				return &Order{
						Amount:    o.Amount - other.Amount,
						Price:     o.Price,
						ID:        o.ID,
						Side:      o.Side,
						User:      o.User,
						Timestamp: o.Timestamp,
					}, &Trade{
						TakeOrderID:  o.ID,
						MakerOrderID: other.ID,
						Amount:       other.Amount,
						Price:        other.Price,
						Timestamp:    time.Now(),
					}
			} else {
				return &Order{
						Amount:    other.Amount - o.Amount,
						Price:     other.Price,
						ID:        other.ID,
						Side:      other.Side,
						User:      other.User,
						Timestamp: other.Timestamp,
					}, &Trade{
						TakeOrderID:  o.ID,
						MakerOrderID: other.ID,
						Amount:       o.Amount,
						Price:        other.Price,
						Timestamp:    time.Now(),
					}
			}
		}
	case Sell:
		if o.Price <= other.Price {
			if o.Amount == other.Amount {
				return nil, &Trade{
					TakeOrderID:  o.ID,
					MakerOrderID: other.ID,
					Amount:       o.Amount,
					Price:        o.Price,
					Timestamp:    time.Now(),
				}
			} else if o.Amount > other.Amount {
				return &Order{
						Amount:    o.Amount - other.Amount,
						Price:     o.Price,
						ID:        o.ID,
						Side:      o.Side,
						User:      o.User,
						Timestamp: o.Timestamp,
					}, &Trade{
						TakeOrderID:  o.ID,
						MakerOrderID: other.ID,
						Amount:       other.Amount,
						Price:        o.Price,
						Timestamp:    time.Now(),
					}
			} else {
				return &Order{
						Amount:    other.Amount - o.Amount,
						Price:     other.Price,
						ID:        other.ID,
						Side:      other.Side,
						User:      other.User,
						Timestamp: other.Timestamp,
					}, &Trade{
						TakeOrderID:  o.ID,
						MakerOrderID: other.ID,
						Amount:       o.Amount,
						Price:        o.Price,
						Timestamp:    time.Now(),
					}
			}
		}
	default:
		panic(fmt.Sprintf("invalid side: %v", o.Side))
	}
	return nil, nil
}
