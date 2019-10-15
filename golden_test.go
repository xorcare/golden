// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golden

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// _goldie is used for as a tool golden, but inside tests.
var _goldie = _golden

func TestMain(m *testing.M) {
	_goldie.flag = _golden.flag
	_golden.flag = nil
	os.Exit(m.Run())
}

func TestAssert(t *testing.T) {
	type args struct {
		test *bufferTB
		got  []byte
	}
	type readFile struct {
		error error
		bytes []byte
	}
	tests := []struct {
		name     string
		args     args
		readFile readFile
		recover  bool
	}{
		{
			name: "success-assert-nil-with-error-not-exist",
			args: args{
				test: new(bufferTB),
				got:  nil,
			},
			readFile: readFile{
				bytes: nil,
				error: os.ErrNotExist,
			},
			recover: false,
		},
		{
			name: "success-assert-data",
			args: args{
				test: new(bufferTB),
				got:  []byte("golden"),
			},
			readFile: readFile{
				bytes: []byte("golden"),
				error: nil,
			},
			recover: false,
		},
		{
			name: "error-reading-file-permission-denied",
			args: args{
				test: new(bufferTB),
				got:  nil,
			},
			readFile: readFile{
				bytes: nil,
				error: os.ErrPermission,
			},
			recover: true,
		},
		{
			name: "failure-assert-data",
			args: args{
				test: new(bufferTB),
				got:  []byte("golden"),
			},
			readFile: readFile{
				bytes: []byte("Z29sZGVu"),
				error: nil,
			},
			recover: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin := _golden
			defer func() { _golden = origin }()

			_golden.readFile = func(filename string) (bytes []byte, e error) {
				t.Logf(`os.ReadFile(%q) `, filename)
				return tt.readFile.bytes, tt.readFile.error
			}
			defer func() {
				if r := recover(); (r == nil) == tt.recover {
					t.Error(r)
				}
				_goldie.SetTest(t).Assert(tt.args.test.Bytes())
			}()
			tt.args.test.name = t.Name()
			Assert(tt.args.test, tt.args.got)
		})
	}
}

func TestRead(t *testing.T) {
	type args struct {
		test *bufferTB
	}
	type readFile struct {
		error error
		bytes []byte
	}
	tests := []struct {
		name     string
		args     args
		want     []byte
		readFile readFile
		recover  bool
	}{
		{
			name: "success-read-data",
			want: []byte("golden"),
			args: args{
				test: new(bufferTB),
			},
			readFile: readFile{
				bytes: []byte("golden"),
				error: nil,
			},
			recover: false,
		},
		{
			name: "success-read-nil",
			want: nil,
			args: args{
				test: new(bufferTB),
			},
			readFile: readFile{
				bytes: nil,
				error: os.ErrNotExist,
			},
			recover: false,
		},
		{
			name: "error-reading-file-permission-denied",
			want: nil,
			args: args{
				test: new(bufferTB),
			},
			readFile: readFile{
				bytes: nil,
				error: os.ErrPermission,
			},
			recover: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin := _golden
			defer func() { _golden = origin }()

			_golden.readFile = func(filename string) (bytes []byte, e error) {
				t.Logf(`os.ReadFile(%q) `, filename)
				_goldie.SetTest(t).SetPrefix("filename").Assert([]byte(filename))
				return tt.readFile.bytes, tt.readFile.error
			}
			defer func() {
				if r := recover(); (r == nil) == tt.recover {
					t.Error(r)
				}
				_goldie.SetTest(t).Assert(tt.args.test.Bytes())
			}()
			tt.args.test.name = t.Name()
			got := Read(tt.args.test)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		test *bufferTB
		do   func(input []byte) (outcome []byte, err error)
	}
	tests := []struct {
		name    string
		args    args
		recover bool
	}{
		{
			name: "run-without-error",
			args: args{
				test: new(bufferTB),
				do: func(input []byte) (outcome []byte, err error) {
					return nil, nil
				},
			},
			recover: false,
		},
		{
			name: "run-with-error",
			args: args{
				test: new(bufferTB),
				do: func(input []byte) (outcome []byte, err error) {
					return nil, os.ErrClosed
				},
			},
			recover: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin := _golden
			defer func() { _golden = origin }()

			_golden.readFile = func(filename string) (bytes []byte, e error) {
				t.Logf(`os.ReadFile(%q)`, filename)
				return nil, nil
			}
			defer func() {
				if r := recover(); (r == nil) == tt.recover {
					t.Error(r)
				}
				_goldie.SetTest(t).Assert(tt.args.test.Bytes())
			}()
			tt.args.test.name = t.Name()
			Run(tt.args.test, tt.args.do)
		})
	}
}

func TestSetTest(t *testing.T) {
	type args struct {
		test TestingTB
	}
	m := new(bufferTB)
	tests := []struct {
		name string
		args args
		want Tool
	}{
		{
			name: "set-test-nil",
			args: args{},
			want: Tool{},
		},
		{
			name: "set-test-mock",
			args: args{m},
			want: Tool{test: m},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin := _golden
			defer func() { _golden = origin }()

			if got := SetTest(tt.args.test); got.test != tt.want.test {
				t.Errorf("SetTest() = %v, want %v", got.test, tt.want.test)
			}
		})
	}
}

func TestTool_Assert(t *testing.T) {
	type args struct {
		got []byte
	}
	type readFile struct {
		error error
		bytes []byte
	}
	tests := []struct {
		name     string
		args     args
		tool     Tool
		test     bufferTB
		readFile readFile
		recover  bool
	}{
		{
			name: "success-assert-nil-with-error-not-exist",
			args: args{
				got: nil,
			},
			readFile: readFile{
				bytes: nil,
				error: os.ErrNotExist,
			},
			recover: false,
		},
		{
			name: "success-assert-data",
			args: args{
				got: []byte("golden"),
			},
			readFile: readFile{
				bytes: []byte("golden"),
				error: nil,
			},
			recover: false,
		},
		{
			name: "error-reading-file-permission-denied",
			args: args{
				got: nil,
			},
			readFile: readFile{
				bytes: nil,
				error: os.ErrPermission,
			},
			recover: true,
		},
		{
			name: "failure-assert-data",
			args: args{
				got: []byte("golden"),
			},
			readFile: readFile{
				bytes: []byte("Z29sZGVu"),
				error: nil,
			},
			recover: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tool.readFile = func(filename string) (bytes []byte, e error) {
				t.Logf(`os.ReadFile(%q) `, filename)
				return tt.readFile.bytes, tt.readFile.error
			}
			defer func() {
				if r := recover(); (r == nil) == tt.recover {
					t.Error(r)
				}
				_goldie.SetTest(t).Assert(tt.test.Bytes())
			}()
			tt.test.name = t.Name()
			tt.tool.SetTest(&tt.test).Assert(tt.args.got)

		})
	}
}

func TestTool_path(t *testing.T) {
	tests := []struct {
		name string
		path string
		tool Tool
	}{
		{
			name: "empty",
			tool: Tool{},
			path: "TestTool_path/empty.golden",
		},
		{
			name: "default",
			tool: _golden,
			path: "testdata/TestTool_path/default.golden",
		},
		{
			name: "path-target-input",
			tool: _golden.SetTarget(Input),
			path: "testdata/TestTool_path/path-target-input.input",
		},
		{
			name: "path-target-golden",
			tool: _golden.SetTarget(Golden),
			path: "testdata/TestTool_path/path-target-golden.golden",
		},
		{
			name: "path-target-input-prefix-gold",
			tool: _golden.SetTarget(Input).SetPrefix("gold"),
			path: "testdata/TestTool_path/path-target-input-prefix-gold.gold.input",
		},
		{
			name: "path-target-golden-prefix-gold",
			tool: _golden.SetTarget(Golden).SetPrefix("gold"),
			path: "testdata/TestTool_path/path-target-golden-prefix-gold.gold.golden",
		},
		{
			name: "path-prefix-with-spaces",
			tool: _golden.SetTarget(Golden).SetPrefix("path prefix with spaces"),
			path: "testdata/TestTool_path/path-prefix-with-spaces.path_prefix_with_spaces.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tool.SetTest(t).path(); got != tt.path {
				t.Fatalf("error want path %q, actual %q", tt.path, got)
			}
		})
	}
}

func TestTool_Read(t *testing.T) {
	type args struct {
		test *bufferTB
		tar  target
	}
	type readFile struct {
		error error
		bytes []byte
	}
	tests := []struct {
		name     string
		tool     Tool
		args     args
		want     []byte
		readFile readFile
		recover  bool
	}{
		{
			name: "success-read-data",
			want: []byte("golden"),
			args: args{
				test: new(bufferTB),
				tar:  Golden,
			},
			readFile: readFile{
				bytes: []byte("golden"),
				error: nil,
			},
			recover: false,
		},
		{
			name: "success-read-nil",
			want: nil,
			args: args{
				test: new(bufferTB),
				tar:  Golden,
			},
			readFile: readFile{
				bytes: nil,
				error: os.ErrNotExist,
			},
			recover: false,
		},
		{
			name: "error-reading-file-permission-denied",
			want: nil,
			args: args{
				test: new(bufferTB),
				tar:  Golden,
			},
			readFile: readFile{
				bytes: nil,
				error: os.ErrPermission,
			},
			recover: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.test.name = t.Name()
			tt.tool.readFile = func(filename string) (bytes []byte, e error) {
				t.Logf(`os.ReadFile(%q) `, filename)
				return tt.readFile.bytes, tt.readFile.error
			}
			defer func() {
				if r := recover(); (r == nil) == tt.recover {
					t.Error(r)
				}
				_goldie.SetTest(t).Assert(tt.args.test.Bytes())
			}()
			got := tt.tool.SetTest(tt.args.test).SetTarget(tt.args.tar).Read()
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTool_Run(t *testing.T) {
	type args struct {
		do func(input []byte) (got []byte, err error)
	}
	tests := []struct {
		name    string
		tool    Tool
		test    bufferTB
		args    args
		recover bool
	}{
		{
			name: "successful-run",
			args: args{
				do: func(input []byte) (got []byte, err error) {
					return nil, nil
				},
			},
			recover: false,
		},
		{
			name: "fatalities-run",
			args: args{
				do: nil,
			},
			recover: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test.name = t.Name()
			tt.tool.readFile = func(filename string) (bytes []byte, e error) {
				t.Logf(`os.ReadFile(%q)`, filename)
				return nil, nil
			}
			defer func() {
				if r := recover(); (r == nil) == tt.recover {
					t.Error(r)
				}
				_goldie.SetTest(t).Assert(tt.test.Bytes())
			}()
			tt.tool.SetTest(&tt.test).Run(tt.args.do)

		})
	}
}

func TestTool_SetTarget(t *testing.T) {
	type args struct {
		tar target
	}
	tests := []struct {
		name string
		tool Tool
		args args
		want Tool
	}{
		{
			name: "set-input-target",
			args: args{
				tar: Input,
			},
			want: Tool{
				target: Input,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tool.SetTarget(tt.args.tar); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tool.SetTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTool_SetTest(t *testing.T) {
	type args struct {
		t TestingTB
	}
	tests := []struct {
		name string
		tool Tool
		args args
		want Tool
	}{
		{
			name: "set-nil",
			args: args{
				t: nil,
			},
			want: Tool{},
		},
		{
			name: "set-test",
			args: args{
				t: t,
			},
			want: Tool{test: t},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tool.SetTest(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tool.SetTest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTool_Update(t *testing.T) {
	type args struct {
		bytes []byte
	}
	type stat struct {
		fileInfo *FakeStat
		error    error
	}
	fal := false
	tru := true
	tests := []struct {
		name string
		tool Tool
		test bufferTB
		args args
		stat stat
	}{
		{
			name: "not-update-with-nil",
			tool: Tool{flag: nil},
			args: args{
				[]byte("golden"),
			},
		},
		{
			name: "not-update-with-false",
			tool: Tool{flag: &fal},
			args: args{
				[]byte("golden"),
			},
		},
		{
			name: "update-with-true",
			tool: Tool{flag: &tru},
			args: args{
				[]byte("golden"),
			},
			stat: stat{
				fileInfo: &FakeStat{isDir: true},
				error:    nil,
			},
		},
	}
	for _, tt := range tests {
		tt.tool.stat = func(name string) (os.FileInfo, error) {
			t.Logf(`os.Stat(%q)`, name)
			if tt.stat.fileInfo != nil {
				tt.stat.fileInfo.name = name
			}
			return tt.stat.fileInfo, tt.stat.error
		}
		tt.tool.mkdirAll = func(path string, perm os.FileMode) error {
			return nil
		}
		tt.tool.writeFile = func(filename string, data []byte, perm os.FileMode) error {
			t.Logf(`os.WriteFile(%q, %q, %d) `, filename, data, perm)
			return nil
		}
		t.Run(tt.name, func(t *testing.T) {
			tt.test.name = t.Name()
			tt.tool.SetTest(&tt.test).Update(tt.args.bytes)
		})
	}
}

func TestTool_write(t *testing.T) {
	type args struct {
		test  *bufferTB
		tar   target
		bytes []byte
	}
	type stat struct {
		fileInfo *FakeStat
		error    error
	}
	tests := []struct {
		name      string
		tool      Tool
		args      args
		writeFile error
		stat      stat
		recover   bool
	}{
		{
			name: "write-nil",
			args: args{
				test:  new(bufferTB),
				tar:   Golden,
				bytes: nil,
			},
			stat: stat{
				error: os.ErrNotExist,
			},
			recover: false,
		},
		{
			name: "write-nil-with-file-exist",
			args: args{
				test:  new(bufferTB),
				tar:   Golden,
				bytes: nil,
			},
			stat: stat{
				fileInfo: new(FakeStat),
				error:    nil,
			},
			recover: false,
		},
		{
			name: "write-empty",
			args: args{
				test:  new(bufferTB),
				tar:   Golden,
				bytes: []byte{},
			},
			stat: stat{
				fileInfo: new(FakeStat),
			},
			recover: false,
		},
		{
			name: "write-bytes",
			args: args{
				test:  new(bufferTB),
				tar:   Golden,
				bytes: []byte("golden"),
			},
			stat: stat{
				fileInfo: &FakeStat{isDir: true},
			},
			recover: false,
		},
		{
			name: "fatality-error",
			args: args{
				test:  new(bufferTB),
				tar:   Golden,
				bytes: []byte("golden"),
			},
			stat: stat{
				error: os.ErrPermission,
			},
			recover: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.test.name = t.Name()
			tt.tool.writeFile = func(filename string, data []byte, perm os.FileMode) error {
				t.Logf(`os.WriteFile(%q, %q, %d) `, filename, data, perm)
				assert.Equal(t, data, tt.args.bytes)
				if data == nil {
					t.Errorf("you cannot write nil values to a file")
				}
				return tt.writeFile
			}
			tt.tool.mkdirAll = func(path string, perm os.FileMode) error {
				t.Logf(`os.MkdirAll(%q, %d) `, path, perm)
				return nil
			}
			tt.tool.remove = func(name string) error {
				t.Logf(`os.Remove(%q)`, name)
				return nil
			}
			tt.tool.stat = func(name string) (os.FileInfo, error) {
				t.Logf(`os.Stat(%q)`, name)
				if tt.stat.fileInfo != nil {
					tt.stat.fileInfo.name = name
				}
				return tt.stat.fileInfo, tt.stat.error
			}
			defer func() {
				if r := recover(); (r == nil) == tt.recover {
					t.Error(r)
				}
				_goldie.SetTest(t).Assert(tt.args.test.Bytes())
			}()
			tt.tool.SetTest(tt.args.test).
				SetTarget(tt.args.tar).
				write(tt.args.bytes)
		})
	}
}

func TestTool_mkdir(t *testing.T) {
	type args struct {
		loc string
	}
	type stat struct {
		fileInfo *FakeStat
		error    error
	}
	tests := []struct {
		name     string
		tool     Tool
		test     bufferTB
		args     args
		stat     stat
		mkdirAll error
		recover  bool
	}{
		{
			name: "fatality-error",
			args: args{
				loc: filepath.Dir(_golden.SetTest(t).path()),
			},
			stat: stat{
				error: os.ErrPermission,
			},
			recover: true,
		},
		{
			name: "error-file-does-not-exist",
			args: args{
				loc: filepath.Dir(_golden.SetTest(t).path()),
			},
			stat: stat{
				error: os.ErrNotExist,
			},
			recover: false,
		},
		{
			name: "error-dir-is-a-file",
			args: args{
				loc: _golden.SetTest(t).path(),
			},
			stat: stat{
				fileInfo: new(FakeStat),
				error:    nil,
			},
			recover: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tool.mkdirAll = func(path string, perm os.FileMode) error {
				t.Logf(`os.MkdirAll(%q, %d) `, path, perm)
				return tt.mkdirAll
			}
			tt.tool.stat = func(name string) (os.FileInfo, error) {
				t.Logf(`os.Stat(%q)`, name)
				if tt.stat.fileInfo != nil {
					tt.stat.fileInfo.name = name
				}
				return tt.stat.fileInfo, tt.stat.error
			}
			defer func() {
				if r := recover(); (r == nil) == tt.recover {
					t.Error(r)
				}
				_goldie.SetTest(t).Assert(tt.test.Bytes())
			}()
			tt.test.name = t.Name()
			tt.tool.SetTest(&tt.test).mkdir(tt.args.loc)
		})
	}
}

func TestTool_ok(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		tool    Tool
		test    bufferTB
		args    args
		recover bool
	}{
		{
			name: "without-error",
			args: args{
				err: nil,
			},
			recover: false,
		},
		{
			name: "with-error",
			args: args{
				err: os.ErrPermission,
			},
			recover: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); (r == nil) == tt.recover {
					t.Error(r)
				}
				_goldie.SetTest(t).Assert(tt.test.Bytes())
			}()
			tt.test.name = t.Name()
			tt.tool.SetTest(&tt.test).ok(tt.args.err)
		})
	}
}

// FakeStat implements os.FileInfo.
type FakeStat struct {
	name     string
	contents string
	mode     os.FileMode
	offset   int
	isDir    bool
}

// os.FileInfo methods.

func (f *FakeStat) Name() string {
	// A bit of a cheat: we only
	// have a basename, so that's
	// also ok for FileInfo.
	return f.name
}

func (f *FakeStat) Size() int64 {
	return int64(len(f.contents))
}

func (f *FakeStat) Mode() os.FileMode {
	return f.mode
}

func (f *FakeStat) ModTime() time.Time {
	return time.Time{}
}

func (f *FakeStat) IsDir() bool {
	return f.isDir
}

func (f *FakeStat) Sys() interface{} {
	return nil
}

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

func Test_rewrite(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "simple case with spaces",
			arg:  "simple case with spaces",
			want: "simple_case_with_spaces",
		},
		{
			name: "simple case with tab",
			arg:  "simple case with\ttab",
			want: "simple_case_with_tab",
		},
		{
			name: "simple case with new line",
			arg:  "simple case with" + "\n" + "new line",
			want: "simple_case_with_new_line",
		},
		{
			name: "incorrect rune(0)",
			arg:  "simple case with" + string(rune(0)),
			want: `simple_case_with\x00`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rewrite(tt.arg); got != tt.want {
				t.Errorf("rewrite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTool_Equal(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name     string
		args     args
		actual   []byte
		expected []byte
		failed   bool
	}{
		{
			name:     "successful nil-nil",
			expected: nil,
			actual:   nil,
			failed:   false,
		},
		{
			name:     "successful []-[]",
			expected: []byte{},
			actual:   []byte{},
			failed:   false,
		},
		{
			name:     "successful golden-golden",
			expected: []byte("golden"),
			actual:   []byte("golden"),
			failed:   false,
		},
		{
			name:     "failure golden-Z29sZGVu",
			expected: []byte("golden"),
			actual:   []byte("Z29sZGVu"),
			failed:   true,
		},
		{
			name:     "failure golden-nil",
			expected: []byte("golden"),
			actual:   nil,
			failed:   true,
		},
		{
			name:     "failure nil-golden",
			expected: nil,
			actual:   []byte("golden"),
			failed:   true,
		},
		{
			name:     "failure []-nil",
			expected: []byte{},
			actual:   nil,
			failed:   true,
		},
		{
			name:     "failure nil-[]",
			expected: nil,
			actual:   []byte{},
			failed:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &bufferTB{name: t.Name()}
			tool := _golden.SetTest(tb)
			tool.mkdirAll = func(path string, perm os.FileMode) error { return nil }
			tool.readFile = helperOSReadFile(t, tt.expected, nil)

			conclusion := tool.Equal(tt.actual)
			conclusion.Fail()
			if conclusion.Failed() {
				assert.Panics(t, func() {
					conclusion.FailNow()
				})
			} else {
				assert.NotPanics(t, func() {
					conclusion.FailNow()
				})
			}
			if assert.Equal(t, tt.failed, conclusion.Failed()) {
				_goldie.SetTest(t).Equal(tb.Bytes()).FailNow()
			}
		})
	}
}

func TestEqual(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name     string
		args     args
		actual   []byte
		expected []byte
		failed   bool
	}{
		{
			name:     "successful nil-nil",
			expected: nil,
			actual:   nil,
			failed:   false,
		},
		{
			name:     "successful []-[]",
			expected: []byte{},
			actual:   []byte{},
			failed:   false,
		},
		{
			name:     "successful golden-golden",
			expected: []byte("golden"),
			actual:   []byte("golden"),
			failed:   false,
		},
		{
			name:     "failure golden-Z29sZGVu",
			expected: []byte("golden"),
			actual:   []byte("Z29sZGVu"),
			failed:   true,
		},
		{
			name:     "failure golden-nil",
			expected: []byte("golden"),
			actual:   nil,
			failed:   true,
		},
		{
			name:     "failure nil-golden",
			expected: nil,
			actual:   []byte("golden"),
			failed:   true,
		},
		{
			name:     "failure []-nil",
			expected: []byte{},
			actual:   nil,
			failed:   true,
		},
		{
			name:     "failure nil-[]",
			expected: nil,
			actual:   []byte{},
			failed:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin := _golden
			defer func() { _golden = origin }()

			tb := &bufferTB{name: t.Name()}
			_golden.mkdirAll = func(path string, perm os.FileMode) error { return nil }
			_golden.readFile = helperOSReadFile(t, tt.expected, nil)

			conclusion := Equal(tb, tt.actual)
			conclusion.Fail()
			if conclusion.Failed() {
				assert.Panics(t, func() {
					conclusion.FailNow()
				})
			} else {
				assert.NotPanics(t, func() {
					conclusion.FailNow()
				})
			}
			if assert.Equal(t, tt.failed, conclusion.Failed()) {
				_goldie.SetTest(t).Equal(tb.Bytes()).FailNow()
			}
		})
	}
}

func helperOSReadFile(t testing.TB, content []byte, err error) func(string) ([]byte, error) {
	bs := make([]byte, len(content))
	if content == nil {
		bs = nil
		if err == nil {
			// The concept of working with golden files: nil == os.ErrNotExist.
			err = os.ErrNotExist
		}
	} else {
		copy(bs, content)
	}
	return func(filename string) ([]byte, error) {
		t.Logf("os.ReadFile(%q)", filename)
		t.Logf("os.ReadFile.bytes:\n %[1]T %#[1]v\n %[2]T %#[2]v", bs, string(bs))
		t.Logf("os.ReadFile.error: %v", err)
		return bs, err
	}
}
