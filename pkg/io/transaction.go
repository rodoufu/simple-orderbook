package io

import "github.com/rodoufu/simple-orderbook/pkg/entity"

type Transaction interface {
	transaction()
}

type NewOrder struct {
	Transaction
	Order entity.Order
}

type CancelOrder struct {
	Transaction
	User    entity.UserID
	OrderID entity.OrderID
}

type FlushAllOrders struct {
	Transaction
}
