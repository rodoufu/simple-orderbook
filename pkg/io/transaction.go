package io

import (
	"github.com/rodoufu/simple-orderbook/pkg/entity"
)

type Transaction interface {
	transaction()
}

type NewOrderTransaction struct {
	Transaction
	Symbol string
	Order  entity.Order
}

type CancelOrderTransaction struct {
	Transaction
	User    entity.UserID
	OrderID entity.OrderID
}

type FlushAllOrdersTransaction struct {
	Transaction
}

type ErrorTransaction struct {
	Transaction
	Err error
}
