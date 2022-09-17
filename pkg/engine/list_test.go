package engine

import (
	"context"
	"github.com/rodoufu/simple-orderbook/pkg/io"
	"reflect"
	"testing"
	"time"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
	"github.com/rodoufu/simple-orderbook/pkg/event"
)

func Test_listEngine_AddOrder(t *testing.T) {
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
				events: make(chan event.Event, 10),
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
				events: make(chan event.Event, 10),
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

func Test_listEngine_CancelOrder(t *testing.T) {
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
				events: make(chan event.Event, 10),
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
				events: make(chan event.Event, 10),
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
				events: make(chan event.Event, 10),
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
				events: make(chan event.Event, 10),
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

func toListEventsOutput(ctx context.Context, events <-chan event.Event) []string {
	var resp []string
	done := ctx.Done()
	for {
		select {
		case <-done:
			return resp
		case evt, ok := <-events:
			if !ok {
				return resp
			}
			if output := evt.Output(); len(output) > 0 {
				resp = append(resp, output)
			}
		}
	}
}

func Test_listEngine_ProcessTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx          context.Context
		transactions []io.Transaction
	}
	tests := []struct {
		name       string
		engine     *listEngine
		args       args
		wantEvents []string
	}{
		//{
		//	name: "scenario 1 balanced book",
		//	engine: &listEngine{
		//		orders:   map[entity.Side][]entity.Order{},
		//		events:   make(chan event.Event, 50),
		//		orderIDs: map[entity.OrderID]entity.Side{},
		//	},
		//	args: args{
		//		ctx: context.Background(),
		//		transactions: []io.Transaction{
		//			io.NewOrderTransaction{
		//				Symbol: "IBM",
		//				Order: entity.Order{
		//					Amount:    100,
		//					Price:     10,
		//					ID:        1,
		//					Side:      entity.Buy,
		//					User:      1,
		//					Timestamp: time.UnixMilli(1),
		//				},
		//			},
		//			io.NewOrderTransaction{
		//				Symbol: "IBM",
		//				Order: entity.Order{
		//					Amount:    100,
		//					Price:     12,
		//					ID:        2,
		//					Side:      entity.Sell,
		//					User:      1,
		//					Timestamp: time.UnixMilli(2),
		//				},
		//			},
		//			io.NewOrderTransaction{
		//				Symbol: "IBM",
		//				Order: entity.Order{
		//					Amount:    100,
		//					Price:     9,
		//					ID:        101,
		//					Side:      entity.Buy,
		//					User:      2,
		//					Timestamp: time.UnixMilli(3),
		//				},
		//			},
		//			io.NewOrderTransaction{
		//				Symbol: "IBM",
		//				Order: entity.Order{
		//					Amount:    100,
		//					Price:     11,
		//					ID:        102,
		//					Side:      entity.Sell,
		//					User:      2,
		//					Timestamp: time.UnixMilli(4),
		//				},
		//			},
		//			io.NewOrderTransaction{
		//				Symbol: "IBM",
		//				Order: entity.Order{
		//					Amount:    100,
		//					Price:     11,
		//					ID:        3,
		//					Side:      entity.Buy,
		//					User:      1,
		//					Timestamp: time.UnixMilli(5),
		//				},
		//			},
		//			io.NewOrderTransaction{
		//				Symbol: "IBM",
		//				Order: entity.Order{
		//					Amount:    100,
		//					Price:     10,
		//					ID:        103,
		//					Side:      entity.Sell,
		//					User:      2,
		//					Timestamp: time.UnixMilli(6),
		//				},
		//			},
		//			io.NewOrderTransaction{
		//				Symbol: "IBM",
		//				Order: entity.Order{
		//					Amount:    100,
		//					Price:     10,
		//					ID:        4,
		//					Side:      entity.Buy,
		//					User:      1,
		//					Timestamp: time.UnixMilli(7),
		//				},
		//			},
		//			io.NewOrderTransaction{
		//				Symbol: "IBM",
		//				Order: entity.Order{
		//					Amount:    100,
		//					Price:     11,
		//					ID:        104,
		//					Side:      entity.Sell,
		//					User:      2,
		//					Timestamp: time.UnixMilli(8),
		//				},
		//			},
		//			io.FlushAllOrdersTransaction{},
		//		},
		//	},
		//	wantEvents: []string{
		//		"A, 1, 1",
		//		"B, B, 10, 100",
		//		"A, 1, 2",
		//		"B, S, 12, 100",
		//		"A, 2, 101",
		//		"A, 2, 102",
		//		"B, S, 11, 100",
		//		"R, 1, 3",
		//		"R, 2, 103",
		//		"A, 1, 4",
		//		"B, B, 10, 200",
		//		"A, 2, 104",
		//		"B, S, 11, 200",
		//	},
		//},
		{
			name: "scenario 8 balanced book, limit buy partial",
			engine: &listEngine{
				orders:   map[entity.Side][]entity.Order{},
				events:   make(chan event.Event, 50),
				orderIDs: map[entity.OrderID]entity.Side{},
			},
			args: args{
				ctx: context.Background(),
				transactions: []io.Transaction{
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     10,
							ID:        1,
							Side:      entity.Buy,
							User:      1,
							Timestamp: time.UnixMilli(1),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     12,
							ID:        2,
							Side:      entity.Sell,
							User:      1,
							Timestamp: time.UnixMilli(2),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     9,
							ID:        101,
							Side:      entity.Buy,
							User:      2,
							Timestamp: time.UnixMilli(3),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     11,
							ID:        102,
							Side:      entity.Sell,
							User:      2,
							Timestamp: time.UnixMilli(4),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    20,
							Price:     11,
							ID:        3,
							Side:      entity.Buy,
							User:      1,
							Timestamp: time.UnixMilli(5),
						},
					},
					io.FlushAllOrdersTransaction{},
				},
			},
			wantEvents: []string{
				"A, 1, 1",
				"B, B, 10, 100",
				"A, 1, 2",
				"B, S, 12, 100",
				"A, 2, 101",
				"A, 2, 102",
				"B, S, 11, 100",
				"A, 1, 3",
				"T, 1, 3, 2, 102, 11, 20",
				"B, S, 11, 80",
			},
		},
		{
			name: "scenario 9 balanced book, cancel best bid and offer",
			engine: &listEngine{
				orders:   map[entity.Side][]entity.Order{},
				events:   make(chan event.Event, 50),
				orderIDs: map[entity.OrderID]entity.Side{},
			},
			args: args{
				ctx: context.Background(),
				transactions: []io.Transaction{
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     10,
							ID:        1,
							Side:      entity.Buy,
							User:      1,
							Timestamp: time.UnixMilli(1),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     12,
							ID:        2,
							Side:      entity.Sell,
							User:      1,
							Timestamp: time.UnixMilli(2),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     9,
							ID:        101,
							Side:      entity.Buy,
							User:      2,
							Timestamp: time.UnixMilli(3),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     11,
							ID:        102,
							Side:      entity.Sell,
							User:      2,
							Timestamp: time.UnixMilli(4),
						},
					},
					io.CancelOrderTransaction{
						User:    1,
						OrderID: 1,
					},
					io.CancelOrderTransaction{
						User:    2,
						OrderID: 102,
					},
					io.FlushAllOrdersTransaction{},
				},
			},
			wantEvents: []string{
				"A, 1, 1",
				"B, B, 10, 100",
				"A, 1, 2",
				"B, S, 12, 100",
				"A, 2, 101",
				"A, 2, 102",
				"B, S, 11, 100",
				"A, 1, 1",
				"B, B, 9, 100",
				"A, 2, 102",
				"B, S, 12, 100",
			},
		},
		{
			name: "scenario 10 balanced book, cancel behind best bid and offer",
			engine: &listEngine{
				orders:   map[entity.Side][]entity.Order{},
				events:   make(chan event.Event, 50),
				orderIDs: map[entity.OrderID]entity.Side{},
			},
			args: args{
				ctx: context.Background(),
				transactions: []io.Transaction{
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     10,
							ID:        1,
							Side:      entity.Buy,
							User:      1,
							Timestamp: time.UnixMilli(1),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     12,
							ID:        2,
							Side:      entity.Sell,
							User:      1,
							Timestamp: time.UnixMilli(2),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     9,
							ID:        101,
							Side:      entity.Buy,
							User:      2,
							Timestamp: time.UnixMilli(3),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     11,
							ID:        102,
							Side:      entity.Sell,
							User:      2,
							Timestamp: time.UnixMilli(4),
						},
					},
					io.CancelOrderTransaction{
						User:    1,
						OrderID: 2,
					},
					io.CancelOrderTransaction{
						User:    2,
						OrderID: 101,
					},
					io.FlushAllOrdersTransaction{},
				},
			},
			wantEvents: []string{
				"A, 1, 1",
				"B, B, 10, 100",
				"A, 1, 2",
				"B, S, 12, 100",
				"A, 2, 101",
				"A, 2, 102",
				"B, S, 11, 100",
				"A, 1, 2",
				"A, 2, 101",
			},
		},
		{
			name: "scenario 11 balanced book, cancel all bids",
			engine: &listEngine{
				orders:   map[entity.Side][]entity.Order{},
				events:   make(chan event.Event, 50),
				orderIDs: map[entity.OrderID]entity.Side{},
			},
			args: args{
				ctx: context.Background(),
				transactions: []io.Transaction{
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     10,
							ID:        1,
							Side:      entity.Buy,
							User:      1,
							Timestamp: time.UnixMilli(1),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     12,
							ID:        2,
							Side:      entity.Sell,
							User:      1,
							Timestamp: time.UnixMilli(2),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     9,
							ID:        101,
							Side:      entity.Buy,
							User:      2,
							Timestamp: time.UnixMilli(3),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     11,
							ID:        102,
							Side:      entity.Sell,
							User:      2,
							Timestamp: time.UnixMilli(4),
						},
					},
					io.CancelOrderTransaction{
						User:    1,
						OrderID: 1,
					},
					io.CancelOrderTransaction{
						User:    2,
						OrderID: 101,
					},
					io.FlushAllOrdersTransaction{},
				},
			},
			wantEvents: []string{
				"A, 1, 1",
				"B, B, 10, 100",
				"A, 1, 2",
				"B, S, 12, 100",
				"A, 2, 101",
				"A, 2, 102",
				"B, S, 11, 100",
				"A, 1, 1",
				"B, B, 9, 100",
				"A, 2, 101",
				"B, B, -, -",
			},
		},
		{
			name: "scenario 12 balanced book, TOB volume changes",
			engine: &listEngine{
				orders:   map[entity.Side][]entity.Order{},
				events:   make(chan event.Event, 50),
				orderIDs: map[entity.OrderID]entity.Side{},
			},
			args: args{
				ctx: context.Background(),
				transactions: []io.Transaction{
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     10,
							ID:        1,
							Side:      entity.Buy,
							User:      1,
							Timestamp: time.UnixMilli(1),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     12,
							ID:        2,
							Side:      entity.Sell,
							User:      1,
							Timestamp: time.UnixMilli(2),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     9,
							ID:        101,
							Side:      entity.Buy,
							User:      2,
							Timestamp: time.UnixMilli(3),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     11,
							ID:        102,
							Side:      entity.Sell,
							User:      2,
							Timestamp: time.UnixMilli(4),
						},
					},
					io.NewOrderTransaction{
						Symbol: "IBM",
						Order: entity.Order{
							Amount:    100,
							Price:     11,
							ID:        103,
							Side:      entity.Sell,
							User:      2,
							Timestamp: time.UnixMilli(5),
						},
					},
					io.CancelOrderTransaction{
						User:    2,
						OrderID: 103,
					},
					io.CancelOrderTransaction{
						User:    2,
						OrderID: 102,
					},
					io.CancelOrderTransaction{
						User:    1,
						OrderID: 2,
					},
					io.FlushAllOrdersTransaction{},
				},
			},
			wantEvents: []string{
				"A, 1, 1",
				"B, B, 10, 100",
				"A, 1, 2",
				"B, S, 12, 100",
				"A, 2, 101",
				"A, 2, 102",
				"B, S, 11, 100",
				"A, 2, 103",
				"B, S, 11, 200",
				"A, 2, 103",
				"B, S, 11, 100",
				"A, 2, 102",
				"B, S, 12, 100",
				"A, 1, 2",
				"B, S, -, -",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for i, transaction := range tt.args.transactions {
				if err := tt.engine.ProcessTransaction(tt.args.ctx, transaction); err != nil {
					t.Errorf("ProcessTransaction(%d) error = %v", i, err)
				}
			}
			tt.engine.Close()
			gotEvents := toListEventsOutput(tt.args.ctx, tt.engine.events)
			if !reflect.DeepEqual(gotEvents, tt.wantEvents) {
				t.Errorf("ProcessTransaction() events: %v, want: %v", gotEvents, tt.wantEvents)
			}
		})
	}
}
