package mth

import (
	"testing"
)

func TestGreaterCommonDivisor(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{args: args{30,20},want: 10},
		{args: args{20,30},want: 10},
		{args: args{45,15},want: 15},
		{args: args{15,45},want: 15},
		{args: args{3,4},want: 1},
		{args: args{0,4},want: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GreaterCommonDivisor(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("GreaterCommonDivisor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSql(t *testing.T){

}
