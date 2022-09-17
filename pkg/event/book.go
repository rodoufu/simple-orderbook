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
	sideUpper := strings.ToUpper(t.Side.String())
	if t.TotalQuantity == 0 {
		return fmt.Sprintf("B, %v, -, -", sideUpper[0:1])
	} else {
		return fmt.Sprintf("B, %v, %v, %v", sideUpper[0:1], t.Price, t.TotalQuantity)
	}
}
