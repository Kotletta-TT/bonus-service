package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLuhnValid(t *testing.T) {
	type args struct {
		order string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid number",
			args: args{
				order: "0235860772",
			},
			want: true,
		},
		{
			name: "Invalid number",
			args: args{
				order: "854725325031",
			},
			want: false,
		},
		{
			name: "Literal string",
			args: args{
				order: "abfpekd;d23",
			},
			want: false,
		},
		{
			name: "Empty string",
			args: args{
				order: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LuhnValid(tt.args.order); got != tt.want {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
