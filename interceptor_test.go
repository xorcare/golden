package golden

import (
	"fmt"
	"testing"
)

func Test_interceptor_Errorf(t *testing.T) {
	t.Run("by-pointer", func(t *testing.T) {
		i := new(interceptor)
		i.Errorf("%s", t.Name())
		if t.Name() != string(*i) {
			t.Fatalf("%T.Errorf() error got = %q, want %q", *i, t.Name(), string(*i))
		}
	})
	t.Run("by-value", func(t *testing.T) {
		i := interceptor("")
		i.Errorf("%s", t.Name())
		if t.Name() != string(i) {
			t.Fatalf("%T.Errorf() error got = %q, want %q", i, t.Name(), string(i))
		}
	})
}

func Test_interceptor_String(t *testing.T) {
	tests := []struct {
		stringer fmt.Stringer
		want     string
	}{
		{stringer: interceptor(""), want: ""},
		{stringer: new(interceptor), want: ""},
		{stringer: interceptor("golden"), want: "golden"},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := tt.stringer.String(); got != tt.want {
				t.Errorf("interceptor.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
