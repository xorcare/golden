package golden

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_target_String(t *testing.T) {
	tests := []struct {
		target target
		want   string
		runner func(assert.TestingT, assert.PanicTestFunc, ...interface{}) bool
	}{
		{
			target: Golden,
			want:   "golden",
			runner: assert.NotPanics,
		},
		{
			target: Input,
			want:   "input",
			runner: assert.NotPanics,
		},
		{
			target: latest,
			want:   "unsupported target: 2",
			runner: assert.Panics,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			tt.runner(t, func() {
				assert.Equal(t, tt.want, tt.target.String())
			})
		})
	}
}
