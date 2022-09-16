package event

import (
	"github.com/rodoufu/simple-orderbook/pkg/entity"
)

type OrderCancelled struct {
	Event
	Order entity.Order
}

type OrderCreated struct {
	Event
	Order entity.Order
}

type OrderUpdated struct {
	Event
	Order entity.Order
}

type OrderFilled struct {
	Event
	Order entity.Order
}
