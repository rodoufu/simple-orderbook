package io

import (
	"fmt"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
)

type Output interface {
	Output() string
}

type CancelOrderOutput struct {
	Output
	User    entity.UserID
	OrderID entity.OrderID
}

func (c *CancelOrderOutput) String() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("A, %v, %v", c.User, c.OrderID)
}
