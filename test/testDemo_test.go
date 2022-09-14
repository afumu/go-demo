package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_add(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test1",
			args: args{1, 2},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, add(tt.args.x, tt.args.y), tt.want, tt.name)
		})
	}
}
