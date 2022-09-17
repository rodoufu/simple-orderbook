package entity

import (
	"reflect"
	"testing"
	"time"
)

func TestOrder_Less(t *testing.T) {
	t.Parallel()
	type args struct {
		other *Order
	}
	time1 := time.Now()
	time2 := time1.Add(time.Second)
	tests := []struct {
		name  string
		order *Order
		args  args
		want  bool
	}{
		{
			name: "buy (10, time1), (20, time2)",
			order: &Order{
				Price:     10,
				Timestamp: time1,
				Side:      Buy,
			},
			args: args{
				other: &Order{
					Price:     20,
					Timestamp: time2,
					Side:      Buy,
				},
			},
			want: true,
		},
		{
			name: "buy (10, time1), (2, time2)",
			order: &Order{
				Price:     10,
				Timestamp: time1,
				Side:      Buy,
			},
			args: args{
				other: &Order{
					Price:     2,
					Timestamp: time2,
					Side:      Buy,
				},
			},
			want: false,
		},
		{
			name: "buy (10, time1), (10, time2)",
			order: &Order{
				Price:     10,
				Timestamp: time1,
				Side:      Buy,
			},
			args: args{
				other: &Order{
					Price:     10,
					Timestamp: time2,
					Side:      Buy,
				},
			},
			want: true,
		},
		{
			name: "buy (10, time2), (10, time1)",
			order: &Order{
				Price:     10,
				Timestamp: time2,
				Side:      Buy,
			},
			args: args{
				other: &Order{
					Price:     10,
					Timestamp: time1,
					Side:      Buy,
				},
			},
			want: false,
		},
		{
			name: "sell (10, time1), (20, time2)",
			order: &Order{
				Price:     10,
				Timestamp: time1,
				Side:      Sell,
			},
			args: args{
				other: &Order{
					Price:     20,
					Timestamp: time2,
					Side:      Sell,
				},
			},
			want: false,
		},
		{
			name: "sell (10, time1), (2, time2)",
			order: &Order{
				Price:     10,
				Timestamp: time1,
				Side:      Sell,
			},
			args: args{
				other: &Order{
					Price:     2,
					Timestamp: time2,
					Side:      Sell,
				},
			},
			want: true,
		},
		{
			name: "sell (10, time1), (10, time2)",
			order: &Order{
				Price:     10,
				Timestamp: time1,
				Side:      Sell,
			},
			args: args{
				other: &Order{
					Price:     10,
					Timestamp: time2,
					Side:      Sell,
				},
			},
			want: true,
		},
		{
			name: "sell (10, time2), (10, time1)",
			order: &Order{
				Price:     10,
				Timestamp: time2,
				Side:      Sell,
			},
			args: args{
				other: &Order{
					Price:     10,
					Timestamp: time1,
					Side:      Sell,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.order.Less(tt.args.other); got != tt.want {
				t.Errorf("Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_Match(t *testing.T) {
	t.Parallel()
	time1 := time.Now()
	time2 := time1.Add(time.Second)
	type args struct {
		other *Order
	}
	tests := []struct {
		name      string
		order     *Order
		args      args
		wantOrder *Order
		wantTrade *Trade
	}{
		{
			name: "empty",
		},
		{
			name: "same side",
			order: &Order{
				Side: Buy,
			},
			args: args{
				other: &Order{
					Side: Buy,
				},
			},
		},
		{
			name: "buy (10, 10) sell (20, 10) no order, no trade",
			order: &Order{
				Side:   Buy,
				Price:  10,
				Amount: 10,
				ID:     1,
				User:   1,
			},
			args: args{
				other: &Order{
					Side:   Sell,
					Price:  20,
					Amount: 10,
					ID:     2,
					User:   2,
				},
			},
		},
		{
			name: "buy (10, 10) sell (10, 10) same price, no order, trade",
			order: &Order{
				Side:   Buy,
				Price:  10,
				Amount: 10,
				ID:     1,
				User:   1,
			},
			args: args{
				other: &Order{
					Side:   Sell,
					Price:  10,
					Amount: 10,
					ID:     2,
					User:   2,
				},
			},
			wantTrade: &Trade{
				TakeOrderID:  1,
				MakerOrderID: 2,
				Amount:       10,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
		{
			name: "buy (10, 10) sell (10, 10) same price, order larger than book, trade",
			order: &Order{
				Side:      Buy,
				Price:     10,
				Amount:    10,
				ID:        1,
				User:      1,
				Timestamp: time1,
			},
			args: args{
				other: &Order{
					Side:   Sell,
					Price:  10,
					Amount: 9,
					ID:     2,
					User:   2,
				},
			},
			wantOrder: &Order{
				Amount:    1,
				Price:     10,
				ID:        1,
				Side:      Buy,
				User:      1,
				Timestamp: time1,
			},
			wantTrade: &Trade{
				TakeOrderID:  1,
				MakerOrderID: 2,
				Amount:       9,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
		{
			name: "buy (10, 10) sell (10, 10) same price, book order larger, trade",
			order: &Order{
				Side:      Buy,
				Price:     10,
				Amount:    10,
				ID:        1,
				User:      1,
				Timestamp: time1,
			},
			args: args{
				other: &Order{
					Side:      Sell,
					Price:     10,
					Amount:    11,
					ID:        2,
					User:      2,
					Timestamp: time2,
				},
			},
			wantOrder: &Order{
				Amount:    1,
				Price:     10,
				ID:        2,
				Side:      Sell,
				User:      2,
				Timestamp: time2,
			},
			wantTrade: &Trade{
				TakeOrderID:  1,
				MakerOrderID: 2,
				Amount:       10,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
		{
			name: "buy (20, 10) sell (10, 10) buy larger price, no order, trade",
			order: &Order{
				Side:   Buy,
				Price:  20,
				Amount: 10,
				ID:     1,
				User:   1,
			},
			args: args{
				other: &Order{
					Side:   Sell,
					Price:  10,
					Amount: 10,
					ID:     2,
					User:   2,
				},
			},
			wantTrade: &Trade{
				TakeOrderID:  1,
				MakerOrderID: 2,
				Amount:       10,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
		{
			name: "buy (20, 10) sell (10, 10) buy larger price, order larger than book, trade",
			order: &Order{
				Side:      Buy,
				Price:     20,
				Amount:    10,
				ID:        1,
				User:      1,
				Timestamp: time1,
			},
			args: args{
				other: &Order{
					Side:   Sell,
					Price:  10,
					Amount: 9,
					ID:     2,
					User:   2,
				},
			},
			wantOrder: &Order{
				Amount:    1,
				Price:     20,
				ID:        1,
				Side:      Buy,
				User:      1,
				Timestamp: time1,
			},
			wantTrade: &Trade{
				TakeOrderID:  1,
				MakerOrderID: 2,
				Amount:       9,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
		{
			name: "buy (20, 10) sell (10, 10) buy larger price, book order larger, trade",
			order: &Order{
				Side:      Buy,
				Price:     20,
				Amount:    10,
				ID:        1,
				User:      1,
				Timestamp: time1,
			},
			args: args{
				other: &Order{
					Side:      Sell,
					Price:     10,
					Amount:    11,
					ID:        2,
					User:      2,
					Timestamp: time2,
				},
			},
			wantOrder: &Order{
				Amount:    1,
				Price:     10,
				ID:        2,
				Side:      Sell,
				User:      2,
				Timestamp: time2,
			},
			wantTrade: &Trade{
				TakeOrderID:  1,
				MakerOrderID: 2,
				Amount:       10,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
		{
			name: "sell (20, 10) buy (10, 10) no order, no trade",
			order: &Order{
				Side:   Sell,
				Price:  20,
				Amount: 10,
				ID:     2,
				User:   2,
			},
			args: args{
				other: &Order{
					Side:   Buy,
					Price:  10,
					Amount: 10,
					ID:     1,
					User:   1,
				},
			},
		},
		{
			name: "sell (10, 10) buy (10, 10) same price, no order, trade",
			order: &Order{
				Side:   Sell,
				Price:  10,
				Amount: 10,
				ID:     2,
				User:   2,
			},
			args: args{
				other: &Order{
					Side:   Buy,
					Price:  10,
					Amount: 10,
					ID:     1,
					User:   1,
				},
			},
			wantTrade: &Trade{
				TakeOrderID:  2,
				MakerOrderID: 1,
				Amount:       10,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
		{
			name: "sell (10, 10) buy (10, 10) same price, order larger than book, trade",
			order: &Order{
				Side:   Sell,
				Price:  10,
				Amount: 9,
				ID:     2,
				User:   2,
			},
			args: args{
				other: &Order{
					Side:      Buy,
					Price:     10,
					Amount:    10,
					ID:        1,
					User:      1,
					Timestamp: time1,
				},
			},
			wantOrder: &Order{
				Amount:    1,
				Price:     10,
				ID:        1,
				Side:      Buy,
				User:      1,
				Timestamp: time1,
			},
			wantTrade: &Trade{
				TakeOrderID:  2,
				MakerOrderID: 1,
				Amount:       9,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
		{
			name: "sell (10, 10) buy (10, 10) same price, book order larger, trade",
			order: &Order{
				Side:      Sell,
				Price:     10,
				Amount:    11,
				ID:        2,
				User:      2,
				Timestamp: time2,
			},
			args: args{
				other: &Order{
					Side:      Buy,
					Price:     10,
					Amount:    10,
					ID:        1,
					User:      1,
					Timestamp: time1,
				},
			},
			wantOrder: &Order{
				Amount:    1,
				Price:     10,
				ID:        2,
				Side:      Sell,
				User:      2,
				Timestamp: time2,
			},
			wantTrade: &Trade{
				TakeOrderID:  2,
				MakerOrderID: 1,
				Amount:       10,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
		{
			name: "sell (10, 10) buy (20, 10) buy larger price, no order, trade",
			order: &Order{
				Side:   Sell,
				Price:  10,
				Amount: 10,
				ID:     2,
				User:   2,
			},
			args: args{
				other: &Order{
					Side:   Buy,
					Price:  20,
					Amount: 10,
					ID:     1,
					User:   1,
				},
			},
			wantTrade: &Trade{
				TakeOrderID:  2,
				MakerOrderID: 1,
				Amount:       10,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
		{
			name: "sell (10, 10) buy (20, 10) buy larger price, order larger than book, trade",
			order: &Order{
				Side:   Sell,
				Price:  10,
				Amount: 9,
				ID:     2,
				User:   2,
			},
			args: args{
				other: &Order{
					Side:      Buy,
					Price:     20,
					Amount:    10,
					ID:        1,
					User:      1,
					Timestamp: time1,
				},
			},
			wantOrder: &Order{
				Amount:    1,
				Price:     20,
				ID:        1,
				Side:      Buy,
				User:      1,
				Timestamp: time1,
			},
			wantTrade: &Trade{
				TakeOrderID:  2,
				MakerOrderID: 1,
				Amount:       9,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
		{
			name: "sell (10, 10) buy (20, 10) buy larger price, book order larger, trade",
			order: &Order{
				Side:      Sell,
				Price:     10,
				Amount:    11,
				ID:        2,
				User:      2,
				Timestamp: time2,
			},
			args: args{
				other: &Order{
					Side:      Buy,
					Price:     20,
					Amount:    10,
					ID:        1,
					User:      1,
					Timestamp: time1,
				},
			},
			wantOrder: &Order{
				Amount:    1,
				Price:     10,
				ID:        2,
				Side:      Sell,
				User:      2,
				Timestamp: time2,
			},
			wantTrade: &Trade{
				TakeOrderID:  2,
				MakerOrderID: 1,
				Amount:       10,
				Price:        10,
				BuyUserID:    1,
				BuyOrderID:   1,
				SellUserID:   2,
				SellOrderID:  2,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotOrder, gotTrade := tt.order.Match(tt.args.other)
			if gotTrade != nil && tt.wantTrade != nil {
				// Hardcoding the timestamp
				tt.wantTrade.Timestamp = gotTrade.Timestamp
			}
			if !reflect.DeepEqual(gotOrder, tt.wantOrder) {
				t.Errorf("Match() gotOrder = %+v, want %+v", gotOrder, tt.wantOrder)
			}
			if !reflect.DeepEqual(gotTrade, tt.wantTrade) {
				t.Errorf("Match() gotTrade = %+v, want %+v", gotTrade, tt.wantTrade)
			}
		})
	}
}

func TestSide_Opposite(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		s    Side
		want Side
	}{
		{
			name: "buy",
			s:    Buy,
			want: Sell,
		},
		{
			name: "sell",
			s:    Sell,
			want: Buy,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.s.Opposite(); got != tt.want {
				t.Errorf("Opposite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSide_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		s    Side
		want string
	}{
		{
			name: "buy",
			s:    Buy,
			want: "buy",
		},
		{
			name: "sell",
			s:    Sell,
			want: "sell",
		},
		{
			name: "invalid side",
			s:    InvalidSide,
			want: "invalid side",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
