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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotOrder, gotTrade := tt.order.Match(tt.args.other)
			if !reflect.DeepEqual(gotOrder, tt.wantOrder) {
				t.Errorf("Match() got = %v, want %v", gotOrder, tt.wantOrder)
			}
			if !reflect.DeepEqual(gotTrade, tt.wantTrade) {
				t.Errorf("Match() got1 = %v, want %v", gotTrade, tt.wantTrade)
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
