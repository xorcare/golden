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

// Prettier интерфейс для пользовательской реализации форматирования значения.
type Prettier interface {
	Prettify(t TestingTB)
}

// Replacer интерфейс для пользовательской реализации замены значения.
type Replacer interface {
	Replace(t TestingTB, actual interface{})
}

// FileInformer сообщает информацию о том в каком месте нужно заменить значение
// в конструкторе снимка.
type FileInformer interface {
	FileLine() int
	FilePath() string
	FuncName() string
}

// InlineReplacer интерфейс который сообщает информацию о том в каком месте
// нужно заменить значение в конструкторе, что позволяет сделать общий инструмент
// для обновления значений через консольный инструмент golden.
type InlineReplacer interface {
	CallerInfo() FileInformer
}

// Sprinter интерфейс для собственной реализации форматирования значения в код go.
type Sprinter interface {
	Sprint(interface{}) string
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
