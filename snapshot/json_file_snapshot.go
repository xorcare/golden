package snapshot

import (
	"github.com/xorcare/golden"
)

type testingHelper interface {
	Helper()
}

type jsonSnapshot struct{}

func (f *jsonSnapshot) Replace(t golden.TestingTB, actual interface{}) {
	if h, ok := t.(testingHelper); ok {
		h.Helper()
	}
	golden.JSONEq(t, actual.(string)).FailNow()
}

func (f *jsonSnapshot) Equal(t golden.TestingTB, actual interface{}) {
	if h, ok := t.(testingHelper); ok {
		h.Helper()
	}

	golden.JSONEq(t, actual.(string)).FailNow()
}

func JSONFile() golden.Snapshot {
	return new(jsonSnapshot)
}
