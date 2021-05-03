package golden

import (
	"os"
)

// Snapshot is intended to indicate a data source for a test, since it can be a
// string, file, or JSON formatted string, all these options are described by
// separate data types that implements the Snapshot interface. For example:
//	want := snapshot.InlineJSON("{}")
//	// or
//	want := snapshot.FileJSON()
// You can implement your own datatype if needed.
type Snapshot interface {
	Equal(t TestingTB, actual interface{})
}

type Prettier interface {
	Prettify(t TestingTB)
}

type Replacer interface {
	Replace(t TestingTB, actual interface{})
}

func SnapshotEq(t TestingTB, expected Snapshot, actual interface{}) {
	if h, ok := t.(testingHelper); ok {
		h.Helper()
	}

	{
		replacer, isReplacer := expected.(Replacer)
		if os.Getenv(updateEnvName) != "" && isReplacer {
			replacer.Replace(t, actual)
			return
		}
	}

	{
		prettier, isPrettier := expected.(Prettier)
		if os.Getenv("GOLDEN_PRETTIFY") != "" && isPrettier {
			prettier.Prettify(t)
			return
		}
	}

	expected.Equal(t, actual)
}
