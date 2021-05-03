package snapshot_test

import (
	"testing"

	"github.com/xorcare/golden"
	"github.com/xorcare/golden/snapshot"
)

func handlerFirst() string {
	return `{"connected":8512,"registered":100,"android":{"connected":40,"registered":50},"aurora":{"connected":45,"registered":50}}`
}

func handlerSecond() string {
	return `{"connected":16,"registered":100}`
}

func TestOnHTTPHandlerWithGolden(t *testing.T) {
	tests := []struct {
		name string
		want golden.Snapshot
	}{
		{
			name: "Table test with inline data",
			want: snapshot.JSONInline(
				`{
				"android": {
					"connected": 40,
					"registered": 50
				},
				"aurora": {
					"connected": 45,
					"registered": 50
				},
				"connected": 8512,
				"registered": 100
			}`,
			),
		},
		{
			name: "Table test with file",
			want: snapshot.JSONFile(),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				golden.SnapshotEq(t, tt.want, handlerFirst())
			},
		)
	}

	t.Run(
		"Extracted test with inline data", func(t *testing.T) {
			want := snapshot.JSONInline(
				`{
			"connected": 16,
			"registered": 100
		}`,
			)
			golden.SnapshotEq(t, want, handlerSecond())
		},
	)

	t.Run(
		"Extracted test with file", func(t *testing.T) {
			want := snapshot.JSONFile()
			golden.SnapshotEq(t, want, handlerSecond())
		},
	)
}

func TestOnHTTPHandler(t *testing.T) {
	tests := []struct {
		name string
		want golden.Snapshot
	}{
		{
			name: "Table test with inline data",
			want: snapshot.JSONInline(
				`{
				"android": {
					"connected": 40,
					"registered": 50
				},
				"aurora": {
					"connected": 45,
					"registered": 50
				},
				"connected": 8512,
				"registered": 100
			}`,
			),
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				golden.SnapshotEq(t, tt.want, handlerFirst())
			},
		)
	}

	t.Run(
		"Extracted test with inline data", func(t *testing.T) {
			want := snapshot.JSONInline(
				`{
			"connected": 16,
			"registered": 100
		}`,
			)
			golden.SnapshotEq(t, want, handlerSecond())
		},
	)
}
