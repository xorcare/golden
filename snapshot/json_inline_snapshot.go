package snapshot

import (
	"github.com/stretchr/testify/require"

	"github.com/xorcare/golden"
)

type jsonInlineSnapshot struct {
	expected string
}

// Equal implements the golden.Snapshot interface.
func (j jsonInlineSnapshot) Equal(t golden.TestingTB, actual interface{}) {
	if h, ok := t.(testingHelper); ok {
		h.Helper()
	}
	require.JSONEq(t, j.expected, actual.(string))
}

func JSONInline(want string) golden.Snapshot {
	return jsonInlineSnapshot{
		expected: want,
	}
}
