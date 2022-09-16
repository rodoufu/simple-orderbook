package entity

import "time"

type Trade struct {
	TakeOrderID  OrderID
	MakerOrderID OrderID
	Amount       uint64
	Price        uint64
	Timestamp    time.Time
}
