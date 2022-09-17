package event

import (
	"fmt"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
)

// OrderCancelled is emitted when an order is successfully canceled.
type OrderCancelled struct {
	Event
	Order entity.Order
}

func (oc *OrderCancelled) Output() string {
	if oc == nil {
		return ""
	}
	return fmt.Sprintf("A, %v, %v", oc.Order.User, oc.Order.ID)
}

// OrderCreated is emitted when an order is successfully added to the book.
type OrderCreated struct {
	Event
	Order entity.Order
}

func (oc *OrderCreated) Output() string {
	if oc == nil {
		return ""
	}
	return fmt.Sprintf("A, %v, %v", oc.Order.User, oc.Order.ID)
}

// OrderUpdated is emitted when an order changes.
type OrderUpdated struct {
	Event
	Order entity.Order
}

func (ou *OrderUpdated) Output() string {
	return ""
}

// OrderFilled is emitted when an order is successfully filled.
type OrderFilled struct {
	Event
	Order entity.Order
	// Full indicates if the order was fully filled.
	Full bool
}

func (of *OrderFilled) Output() string {
	return ""
}
