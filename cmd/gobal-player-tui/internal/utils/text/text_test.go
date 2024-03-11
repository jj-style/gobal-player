package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatHelp(t *testing.T) {
	type args struct {
		help map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "format shortcuts dictionary",
			args: args{
				help: map[string]string{"a": "hello", "b": "world"},
			},
			want: "a: hello, b: world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatHelp(tt.args.help)
			assert.Equal(t, tt.want, got)
		})
	}
}
