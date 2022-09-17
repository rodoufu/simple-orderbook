package orderbook

import (
	"context"
	"github.com/rodoufu/simple-orderbook/pkg/entity"
	"reflect"
	"sync"
	"testing"
)

func toListBookLevel(ctx context.Context, books <-chan BookLevel) []BookLevel {
	var resp []BookLevel
	done := ctx.Done()
	for {
		select {
		case <-done:
			return resp
		case book, ok := <-books:
			if !ok {
				return resp
			}
			resp = append(resp, book)
		}
	}
}

func Test_listOrderBook_Bids(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		orderBook *listOrderBook
		args      args
		want      []BookLevel
	}{
		{
			name: "empty",
			orderBook: &listOrderBook{
				mtx: map[entity.Side]*sync.RWMutex{
					entity.Buy:  {},
					entity.Sell: {},
				},
			},
			args: args{ctx: context.Background()},
		},
		{
			name: "two orders",
			orderBook: &listOrderBook{
				mtx: map[entity.Side]*sync.RWMutex{
					entity.Buy:  {},
					entity.Sell: {},
				},
				orders: map[entity.Side][]entity.Order{
					entity.Buy: {
						{
							Amount: 10,
							Price:  10,
							ID:     1,
							Side:   entity.Buy,
							User:   1,
						},
						{
							Amount: 10,
							Price:  10,
							ID:     2,
							Side:   entity.Buy,
							User:   2,
						},
					},
				},
			},
			args: args{ctx: context.Background()},
			want: []BookLevel{
				{
					Side:          entity.Buy,
					Price:         10,
					TotalQuantity: 20,
				},
			},
		},
		{
			name: "3 orders",
			orderBook: &listOrderBook{
				mtx: map[entity.Side]*sync.RWMutex{
					entity.Buy:  {},
					entity.Sell: {},
				},
				orders: map[entity.Side][]entity.Order{
					entity.Buy: {
						{
							Amount: 10,
							Price:  10,
							ID:     1,
							Side:   entity.Buy,
							User:   1,
						},
						{
							Amount: 10,
							Price:  10,
							ID:     2,
							Side:   entity.Buy,
							User:   2,
						},
						{
							Amount: 10,
							Price:  11,
							ID:     3,
							Side:   entity.Buy,
							User:   3,
						},
					},
				},
			},
			args: args{ctx: context.Background()},
			want: []BookLevel{
				{
					Side:          entity.Buy,
					Price:         11,
					TotalQuantity: 10,
				},
				{
					Side:          entity.Buy,
					Price:         10,
					TotalQuantity: 20,
				},
			},
		},
		{
			name: "4 orders",
			orderBook: &listOrderBook{
				mtx: map[entity.Side]*sync.RWMutex{
					entity.Buy:  {},
					entity.Sell: {},
				},
				orders: map[entity.Side][]entity.Order{
					entity.Buy: {
						{
							Amount: 10,
							Price:  10,
							ID:     1,
							Side:   entity.Buy,
							User:   1,
						},
						{
							Amount: 10,
							Price:  10,
							ID:     2,
							Side:   entity.Buy,
							User:   2,
						},
						{
							Amount: 10,
							Price:  11,
							ID:     3,
							Side:   entity.Buy,
							User:   3,
						},
						{
							Amount: 15,
							Price:  11,
							ID:     4,
							Side:   entity.Buy,
							User:   4,
						},
					},
				},
			},
			args: args{ctx: context.Background()},
			want: []BookLevel{
				{
					Side:          entity.Buy,
					Price:         11,
					TotalQuantity: 25,
				},
				{
					Side:          entity.Buy,
					Price:         10,
					TotalQuantity: 20,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := toListBookLevel(tt.args.ctx, tt.orderBook.Bids(tt.args.ctx)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bids() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
