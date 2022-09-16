package entity

import (
	"reflect"
	"testing"
)

func TestOrder_Less(t *testing.T) {
	t.Parallel()
	type args struct {
		other *Order
	}
	tests := []struct {
		name  string
		order *Order
		args  args
		want  bool
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
