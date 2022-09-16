package engine

import (
	"context"
	"reflect"
	"testing"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
	"github.com/rodoufu/simple-orderbook/pkg/event"
)

func Test_simpleEngine_AddOrder(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx   context.Context
		order entity.Order
	}
	tests := []struct {
		name       string
		engine     *listEngine
		args       args
		wantErr    bool
		wantOrders map[entity.Side][]entity.Order
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name: "empty book, add sell",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {},
					entity.Buy:  {},
				},
				events: make(chan event.Event, 1),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 10,
					Price:  100,
					ID:     1,
					Side:   entity.Sell,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						Amount: 10,
						Price:  100,
						ID:     1,
						Side:   entity.Sell,
						User:   1,
					},
				},
				entity.Buy: {},
			},
		},
		{
			name: "empty book, add buy",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {},
					entity.Buy:  {},
				},
				events: make(chan event.Event, 1),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 10,
					Price:  100,
					ID:     1,
					Side:   entity.Buy,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Buy: {
					{
						Amount: 10,
						Price:  100,
						ID:     1,
						Side:   entity.Buy,
						User:   1,
					},
				},
				entity.Sell: {},
			},
		},
		{
			name: "add sell",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							Amount: 9,
							Price:  90,
							ID:     2,
							Side:   entity.Sell,
							User:   2,
						},
					},
					entity.Buy: {},
				},
				events: make(chan event.Event, 1),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 10,
					Price:  100,
					ID:     1,
					Side:   entity.Sell,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						Amount: 10,
						Price:  100,
						ID:     1,
						Side:   entity.Sell,
						User:   1,
					},
					{
						Amount: 9,
						Price:  90,
						ID:     2,
						Side:   entity.Sell,
						User:   2,
					},
				},
				entity.Buy: {},
			},
		},
		{
			name: "add buy",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {},
					entity.Buy: {
						{
							Amount: 9,
							Price:  90,
							ID:     2,
							Side:   entity.Buy,
							User:   2,
						},
					},
				},
				events: make(chan event.Event, 1),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 10,
					Price:  100,
					ID:     1,
					Side:   entity.Buy,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Buy: {
					{
						Amount: 9,
						Price:  90,
						ID:     2,
						Side:   entity.Buy,
						User:   2,
					},
					{
						Amount: 10,
						Price:  100,
						ID:     1,
						Side:   entity.Buy,
						User:   1,
					},
				},
				entity.Sell: {},
			},
		},
		{
			name: "add buy, match, full fill",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							Amount: 10,
							Price:  200,
							ID:     2,
							Side:   entity.Sell,
							User:   2,
						},
						{
							Amount: 9,
							Price:  150,
							ID:     3,
							Side:   entity.Sell,
							User:   3,
						},
					},
					entity.Buy: {
						{
							Amount: 9,
							Price:  90,
							ID:     4,
							Side:   entity.Buy,
							User:   4,
						},
						{
							Amount: 10,
							Price:  100,
							ID:     5,
							Side:   entity.Buy,
							User:   5,
						},
					},
				},
				events: make(chan event.Event, 10),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 9,
					Price:  150,
					ID:     1,
					Side:   entity.Buy,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						Amount: 10,
						Price:  200,
						ID:     2,
						Side:   entity.Sell,
						User:   2,
					},
				},
				entity.Buy: {
					{
						Amount: 9,
						Price:  90,
						ID:     4,
						Side:   entity.Buy,
						User:   4,
					},
					{
						Amount: 10,
						Price:  100,
						ID:     5,
						Side:   entity.Buy,
						User:   5,
					},
				},
			},
		},
		{
			name: "add sell, match, full fill",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							Amount: 10,
							Price:  200,
							ID:     2,
							Side:   entity.Sell,
							User:   2,
						},
						{
							Amount: 9,
							Price:  150,
							ID:     3,
							Side:   entity.Sell,
							User:   3,
						},
					},
					entity.Buy: {
						{
							Amount: 9,
							Price:  90,
							ID:     4,
							Side:   entity.Buy,
							User:   4,
						},
						{
							Amount: 10,
							Price:  100,
							ID:     5,
							Side:   entity.Buy,
							User:   5,
						},
					},
				},
				events: make(chan event.Event, 10),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 10,
					Price:  100,
					ID:     1,
					Side:   entity.Sell,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						Amount: 10,
						Price:  200,
						ID:     2,
						Side:   entity.Sell,
						User:   2,
					},
					{
						Amount: 9,
						Price:  150,
						ID:     3,
						Side:   entity.Sell,
						User:   3,
					},
				},
				entity.Buy: {
					{
						Amount: 9,
						Price:  90,
						ID:     4,
						Side:   entity.Buy,
						User:   4,
					},
				},
			},
		},
		{
			name: "add buy, match, full fill, book is larger",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							Amount: 10,
							Price:  200,
							ID:     2,
							Side:   entity.Sell,
							User:   2,
						},
						{
							Amount: 10,
							Price:  150,
							ID:     3,
							Side:   entity.Sell,
							User:   3,
						},
					},
					entity.Buy: {
						{
							Amount: 9,
							Price:  90,
							ID:     4,
							Side:   entity.Buy,
							User:   4,
						},
						{
							Amount: 10,
							Price:  100,
							ID:     5,
							Side:   entity.Buy,
							User:   5,
						},
					},
				},
				events: make(chan event.Event, 10),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 9,
					Price:  150,
					ID:     1,
					Side:   entity.Buy,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						Amount: 10,
						Price:  200,
						ID:     2,
						Side:   entity.Sell,
						User:   2,
					},
					{
						Amount: 1,
						Price:  150,
						ID:     3,
						Side:   entity.Sell,
						User:   3,
					},
				},
				entity.Buy: {
					{
						Amount: 9,
						Price:  90,
						ID:     4,
						Side:   entity.Buy,
						User:   4,
					},
					{
						Amount: 10,
						Price:  100,
						ID:     5,
						Side:   entity.Buy,
						User:   5,
					},
				},
			},
		},
		{
			name: "add sell, match, full fill, book is larger",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							Amount: 10,
							Price:  200,
							ID:     2,
							Side:   entity.Sell,
							User:   2,
						},
						{
							Amount: 9,
							Price:  150,
							ID:     3,
							Side:   entity.Sell,
							User:   3,
						},
					},
					entity.Buy: {
						{
							Amount: 9,
							Price:  90,
							ID:     4,
							Side:   entity.Buy,
							User:   4,
						},
						{
							Amount: 11,
							Price:  100,
							ID:     5,
							Side:   entity.Buy,
							User:   5,
						},
					},
				},
				events: make(chan event.Event, 10),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 10,
					Price:  100,
					ID:     1,
					Side:   entity.Sell,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						Amount: 10,
						Price:  200,
						ID:     2,
						Side:   entity.Sell,
						User:   2,
					},
					{
						Amount: 9,
						Price:  150,
						ID:     3,
						Side:   entity.Sell,
						User:   3,
					},
				},
				entity.Buy: {
					{
						Amount: 9,
						Price:  90,
						ID:     4,
						Side:   entity.Buy,
						User:   4,
					},
					{
						Amount: 1,
						Price:  100,
						ID:     5,
						Side:   entity.Buy,
						User:   5,
					},
				},
			},
		},
		{
			name: "add buy, match, full fill, order is larger",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							Amount: 10,
							Price:  200,
							ID:     2,
							Side:   entity.Sell,
							User:   2,
						},
						{
							Amount: 9,
							Price:  150,
							ID:     3,
							Side:   entity.Sell,
							User:   3,
						},
					},
					entity.Buy: {
						{
							Amount: 9,
							Price:  90,
							ID:     4,
							Side:   entity.Buy,
							User:   4,
						},
						{
							Amount: 10,
							Price:  100,
							ID:     5,
							Side:   entity.Buy,
							User:   5,
						},
					},
				},
				events: make(chan event.Event, 10),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 10,
					Price:  150,
					ID:     1,
					Side:   entity.Buy,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						Amount: 10,
						Price:  200,
						ID:     2,
						Side:   entity.Sell,
						User:   2,
					},
				},
				entity.Buy: {
					{
						Amount: 9,
						Price:  90,
						ID:     4,
						Side:   entity.Buy,
						User:   4,
					},
					{
						Amount: 10,
						Price:  100,
						ID:     5,
						Side:   entity.Buy,
						User:   5,
					},
					{
						Amount: 1,
						Price:  150,
						ID:     1,
						Side:   entity.Buy,
						User:   1,
					},
				},
			},
		},
		{
			name: "add sell, match, full fill, order is larger",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							Amount: 10,
							Price:  200,
							ID:     2,
							Side:   entity.Sell,
							User:   2,
						},
						{
							Amount: 9,
							Price:  150,
							ID:     3,
							Side:   entity.Sell,
							User:   3,
						},
					},
					entity.Buy: {
						{
							Amount: 9,
							Price:  90,
							ID:     4,
							Side:   entity.Buy,
							User:   4,
						},
						{
							Amount: 11,
							Price:  100,
							ID:     5,
							Side:   entity.Buy,
							User:   5,
						},
					},
				},
				events: make(chan event.Event, 10),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 12,
					Price:  100,
					ID:     1,
					Side:   entity.Sell,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						Amount: 10,
						Price:  200,
						ID:     2,
						Side:   entity.Sell,
						User:   2,
					},
					{
						Amount: 9,
						Price:  150,
						ID:     3,
						Side:   entity.Sell,
						User:   3,
					},
					{
						Amount: 1,
						Price:  100,
						ID:     1,
						Side:   entity.Sell,
						User:   1,
					},
				},
				entity.Buy: {
					{
						Amount: 9,
						Price:  90,
						ID:     4,
						Side:   entity.Buy,
						User:   4,
					},
				},
			},
		},
		{
			name: "add buy, match, full fill two from the book, order is larger",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							Amount: 10,
							Price:  200,
							ID:     2,
							Side:   entity.Sell,
							User:   2,
						},
						{
							Amount: 9,
							Price:  150,
							ID:     3,
							Side:   entity.Sell,
							User:   3,
						},
						{
							Amount: 11,
							Price:  120,
							ID:     3,
							Side:   entity.Sell,
							User:   3,
						},
					},
					entity.Buy: {
						{
							Amount: 9,
							Price:  90,
							ID:     4,
							Side:   entity.Buy,
							User:   4,
						},
						{
							Amount: 10,
							Price:  100,
							ID:     5,
							Side:   entity.Buy,
							User:   5,
						},
					},
				},
				events: make(chan event.Event, 10),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 21,
					Price:  150,
					ID:     1,
					Side:   entity.Buy,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						Amount: 10,
						Price:  200,
						ID:     2,
						Side:   entity.Sell,
						User:   2,
					},
				},
				entity.Buy: {
					{
						Amount: 9,
						Price:  90,
						ID:     4,
						Side:   entity.Buy,
						User:   4,
					},
					{
						Amount: 10,
						Price:  100,
						ID:     5,
						Side:   entity.Buy,
						User:   5,
					},
					{
						Amount: 1,
						Price:  150,
						ID:     1,
						Side:   entity.Buy,
						User:   1,
					},
				},
			},
		},
		{
			name: "add buy, match, full fill 3 from the book, order is larger",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							Amount: 10,
							Price:  200,
							ID:     2,
							Side:   entity.Sell,
							User:   2,
						},
						{
							Amount: 9,
							Price:  150,
							ID:     3,
							Side:   entity.Sell,
							User:   3,
						},
						{
							Amount: 10,
							Price:  130,
							ID:     3,
							Side:   entity.Sell,
							User:   3,
						},
						{
							Amount: 11,
							Price:  120,
							ID:     3,
							Side:   entity.Sell,
							User:   3,
						},
					},
					entity.Buy: {
						{
							Amount: 9,
							Price:  90,
							ID:     4,
							Side:   entity.Buy,
							User:   4,
						},
						{
							Amount: 10,
							Price:  100,
							ID:     5,
							Side:   entity.Buy,
							User:   5,
						},
					},
				},
				events: make(chan event.Event, 10),
			},
			args: args{
				ctx: context.Background(),
				order: entity.Order{
					Amount: 31,
					Price:  150,
					ID:     1,
					Side:   entity.Buy,
					User:   1,
				},
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						Amount: 10,
						Price:  200,
						ID:     2,
						Side:   entity.Sell,
						User:   2,
					},
				},
				entity.Buy: {
					{
						Amount: 9,
						Price:  90,
						ID:     4,
						Side:   entity.Buy,
						User:   4,
					},
					{
						Amount: 10,
						Price:  100,
						ID:     5,
						Side:   entity.Buy,
						User:   5,
					},
					{
						Amount: 1,
						Price:  150,
						ID:     1,
						Side:   entity.Buy,
						User:   1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.engine.AddOrder(tt.args.ctx, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("AddOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.engine != nil && !reflect.DeepEqual(tt.wantOrders, tt.engine.orders) {
				t.Errorf("CancelOrder() got = %+v, want %+v", tt.engine.orders, tt.wantOrders)
			}
		})
	}
}

func Test_simpleEngine_CancelOrder(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context
		orderID entity.OrderID
	}
	tests := []struct {
		name       string
		engine     *listEngine
		args       args
		wantErr    bool
		wantOrders map[entity.Side][]entity.Order
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name:   "order does not exist",
			engine: &listEngine{},
			args: args{
				ctx:     context.Background(),
				orderID: 1,
			},
			wantErr: true,
		},
		{
			name: "delete sell order - get empty book",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{
					1: entity.Sell,
				},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							ID: 1,
						},
					},
				},
				events: make(chan event.Event, 1),
			},
			args: args{
				ctx:     context.Background(),
				orderID: 1,
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {},
			},
		},
		{
			name: "delete sell order - last order",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{
					1: entity.Sell,
					2: entity.Sell,
				},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							ID: 2,
						},
						{
							ID: 1,
						},
					},
				},
				events: make(chan event.Event, 1),
			},
			args: args{
				ctx:     context.Background(),
				orderID: 1,
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						ID: 2,
					},
				},
			},
		},
		{
			name: "delete sell order - first order",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{
					1: entity.Sell,
					2: entity.Sell,
				},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							ID: 2,
						},
						{
							ID: 1,
						},
					},
				},
				events: make(chan event.Event, 1),
			},
			args: args{
				ctx:     context.Background(),
				orderID: 2,
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						ID: 1,
					},
				},
			},
		},
		{
			name: "delete sell order - middle order",
			engine: &listEngine{
				orderIDs: map[entity.OrderID]entity.Side{
					1: entity.Sell,
					2: entity.Sell,
					3: entity.Sell,
				},
				orders: map[entity.Side][]entity.Order{
					entity.Sell: {
						{
							ID: 3,
						},
						{
							ID: 2,
						},
						{
							ID: 1,
						},
					},
				},
				events: make(chan event.Event, 1),
			},
			args: args{
				ctx:     context.Background(),
				orderID: 2,
			},
			wantErr: false,
			wantOrders: map[entity.Side][]entity.Order{
				entity.Sell: {
					{
						ID: 3,
					},
					{
						ID: 1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.engine.CancelOrder(tt.args.ctx, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("CancelOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.engine != nil && !reflect.DeepEqual(tt.wantOrders, tt.engine.orders) {
				t.Errorf("CancelOrder() got = %+v, want %+v", tt.engine.orders, tt.wantOrders)
			}
		})
	}
}
