package golden

import (
	"reflect"
	"testing"
)

func Test_target_String(t *testing.T) {
	tests := []struct {
		name    string
		t       target
		want    string
		recover bool
	}{
		{
			name:    "Golden",
			t:       Golden,
			want:    "golden",
			recover: false,
		},
		{
			name:    "Input",
			t:       Input,
			want:    "input",
			recover: false,
		},
		{
			name:    "Panic",
			t:       latest,
			want:    "unsupported target: 2",
			recover: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					} else if r != nil && !reflect.DeepEqual(r, tt.want) {
						t.Errorf("the expected result of execution = %v, want %v", r, tt.want)
					}
				}()
				if got := tt.t.String(); got != tt.want {
					t.Errorf("target.String() = %v, want %v", got, tt.want)
				}
			}()
		})
	}
}
