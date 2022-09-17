package event

import (
	"github.com/rodoufu/simple-orderbook/pkg/entity"
)

// OrderCancelled is emitted when an order is successfully canceled.
type OrderCancelled struct {
	Event
	Order entity.Order
}

// OrderCreated is emitted when an order is successfully added to the book.
type OrderCreated struct {
	Event
	Order entity.Order
}

// OrderUpdated is emitted when an order changes.
type OrderUpdated struct {
	Event
	Order entity.Order
}

// OrderFilled is emitted when an order is successfully filled.
type OrderFilled struct {
	Event
	Order entity.Order
	// Full indicates if the order was fully filled.
	Full bool
}
