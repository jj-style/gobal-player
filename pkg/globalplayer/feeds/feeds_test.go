package feeds

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFmtDuration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		given time.Duration
		want  string
	}{
		{name: "seconds", given: time.Second * 3, want: "00:00:03"},
		{name: "minutes and seconds", given: time.Minute + (time.Second * 15), want: "00:01:15"},
		{name: "hours minutes and seconds", given: (time.Hour * 2) + (time.Minute * 30) + (time.Second * 5), want: "02:30:05"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, fmtDuration(tt.given))
		})
	}
}
