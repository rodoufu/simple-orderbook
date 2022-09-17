package event

import (
	"fmt"
	"strings"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
)

type TopOfBookChange struct {
	Event
	Side          entity.Side
	Price         uint64
	TotalQuantity uint64
}

func (t *TopOfBookChange) Output() string {
	if t == nil {
		return ""
	}
	if t.TotalQuantity == 0 {
		return fmt.Sprintf("B, %v, -, -", strings.ToUpper(t.Side.String())[1])
	} else {
		return fmt.Sprintf("B, %v, %v, %v", strings.ToUpper(t.Side.String())[1], t.Price, t.TotalQuantity)
	}
}
