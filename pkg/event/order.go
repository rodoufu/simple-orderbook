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
	return ""
}

// OrderCreated is emitted when an order is successfully added to the book.
type OrderCreated struct {
	Event
	Order entity.Order
}

func (oc *OrderCreated) Output() string {
	return ""
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

// OrderAcknowledge is used only to print messages.
type OrderAcknowledge struct {
	Event
	Order entity.Order
}

func (oa *OrderAcknowledge) Output() string {
	if oa == nil {
		return ""
	}
	return fmt.Sprintf("A, %v, %v", oa.Order.User, oa.Order.ID)
}
