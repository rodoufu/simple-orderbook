package entity

import "time"

// Trade represents a trade generated on a match.
type Trade struct {
	// TakeOrderID is the order being added to the book.
	TakeOrderID OrderID
	// MakerOrderID is the order already on the book.
	MakerOrderID OrderID
	// Amount is the size of the trade.
	Amount uint64
	// Price is how much the client paid for the trade.
	Price uint64
	// Timestamp is the moment the trade was created.
	Timestamp   time.Time
	BuyUserID   UserID
	BuyOrderID  OrderID
	SellUserID  UserID
	SellOrderID OrderID
}
