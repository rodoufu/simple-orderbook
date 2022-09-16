package entity

import (
	"fmt"
	"time"
)

type Side uint8

const (
	InvalidSide Side = iota
	Buy         Side = iota
	Sell        Side = iota
)

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

type OrderID uint64
type UserID uint64

type Order struct {
	Amount    uint64
	Price     uint64
	ID        OrderID
	Side      Side
	User      UserID
	Timestamp time.Time
}

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
