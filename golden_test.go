// Copyright (c) 2019-2021 Vasiliy Vasilyuk. All rights reserved.
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
		path string
		tool Tool
	}{
		{
			tool: Tool{},
			path: "TestTool_path/#00.golden",
		},
		{
			tool: _golden,
			path: "testdata/TestTool_path/#01.golden",
		},
		{
			tool: _golden.SetTarget(Input),
			path: "testdata/TestTool_path/#02.input",
		},
		{
			tool: _golden.SetTarget(Golden),
			path: "testdata/TestTool_path/#03.golden",
		},
		{
			tool: _golden.SetTarget(Input).SetPrefix("prefix"),
			path: "testdata/TestTool_path/#04.prefix.input",
		},
		{
			tool: _golden.SetTarget(Golden).SetPrefix("prefix"),
			path: "testdata/TestTool_path/#05.prefix.golden",
		},
		{
			tool: _golden.SetTarget(Golden).SetPrefix("prefix with spaces"),
			path: "testdata/TestTool_path/#06.prefix_with_spaces.golden",
		},
		{
			tool: _golden.setExtension("extension").SetTarget(Input).SetPrefix("prefix"),
			path: "testdata/TestTool_path/#07.prefix.extension.input",
		},
		{
			tool: _golden.setExtension("extension").SetTarget(Golden).SetPrefix("prefix"),
			path: "testdata/TestTool_path/#08.prefix.extension.golden",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, tt.path, tt.tool.SetTest(t).path())
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
	t.Run("with-set-want-field", func(t *testing.T) {
		tb := &bufferTB{name: t.Name()}
		tool := SetWant(tb, []byte(t.Name()))
		assert.Equal(t, t.Name(), string(tool.Read()))
		_goldie.SetTest(t).Assert(tb.Bytes())
	})
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
	tests := []struct {
		name   string
		tool   Tool
		target target
		want   Tool
	}{
		{
			name:   "set-input-target",
			target: Input,
			want:   Tool{target: Input},
		},
		{
			name:   "set-golden-target",
			target: Golden,
			want:   Tool{target: Golden},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.tool.SetTarget(tt.target))
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

func TestTool_noError(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		runner func(assert.TestingT, assert.PanicTestFunc, ...interface{}) bool
	}{
		{
			name:   "without-error",
			err:    nil,
			runner: assert.NotPanics,
		},
		{
			name:   "with-error",
			err:    os.ErrPermission,
			runner: assert.Panics,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &bufferTB{name: t.Name()}
			tt.runner(t, func() {
				SetTest(tb).noError(tt.err)
			})
			_goldie.SetTest(t).Equal(tb.Bytes()).FailNow()
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
	// also noError for FileInfo.
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
		name   string
		args   args
		got    []byte
		want   []byte
		failed bool
	}{
		{
			name:   "successful nil-nil",
			want:   nil,
			got:    nil,
			failed: false,
		},
		{
			name:   "successful []-[]",
			want:   []byte{},
			got:    []byte{},
			failed: false,
		},
		{
			name:   "successful golden-golden",
			want:   []byte("golden"),
			got:    []byte("golden"),
			failed: false,
		},
		{
			name:   "failure golden-Z29sZGVu",
			want:   []byte("golden"),
			got:    []byte("Z29sZGVu"),
			failed: true,
		},
		{
			name:   "failure golden-nil",
			want:   []byte("golden"),
			got:    nil,
			failed: true,
		},
		{
			name:   "failure nil-golden",
			want:   nil,
			got:    []byte("golden"),
			failed: true,
		},
		{
			name:   "failure []-nil",
			want:   []byte{},
			got:    nil,
			failed: true,
		},
		{
			name:   "failure nil-[]",
			want:   nil,
			got:    []byte{},
			failed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &bufferTB{name: t.Name()}
			tool := _golden.SetTest(tb)
			tool.mkdirAll = func(path string, perm os.FileMode) error { return nil }
			tool.readFile = helperOSReadFile(t, tt.want, nil)

			conclusion := tool.Equal(tt.got)
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
		name   string
		args   args
		got    []byte
		want   []byte
		failed bool
	}{
		{
			name:   "successful nil-nil",
			want:   nil,
			got:    nil,
			failed: false,
		},
		{
			name:   "successful []-[]",
			want:   []byte{},
			got:    []byte{},
			failed: false,
		},
		{
			name:   "successful golden-golden",
			want:   []byte("golden"),
			got:    []byte("golden"),
			failed: false,
		},
		{
			name:   "failure golden-Z29sZGVu",
			want:   []byte("golden"),
			got:    []byte("Z29sZGVu"),
			failed: true,
		},
		{
			name:   "failure golden-nil",
			want:   []byte("golden"),
			got:    nil,
			failed: true,
		},
		{
			name:   "failure nil-golden",
			want:   nil,
			got:    []byte("golden"),
			failed: true,
		},
		{
			name:   "failure []-nil",
			want:   []byte{},
			got:    nil,
			failed: true,
		},
		{
			name:   "failure nil-[]",
			want:   nil,
			got:    []byte{},
			failed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin := _golden
			defer func() { _golden = origin }()

			tb := &bufferTB{name: t.Name()}
			_golden.mkdirAll = func(path string, perm os.FileMode) error { return nil }
			_golden.readFile = helperOSReadFile(t, tt.want, nil)

			conclusion := Equal(tb, tt.got)
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

func Test_jsonFormatter(t *testing.T) {
	tests := []struct {
		name string
		json string
		want string
	}{
		{
			json: `{}`,
			want: `{}`,
		},
		{
			json: `{"data":null}`,
			want: "{\n\t\"data\": null\n}",
		},
		{
			json: `{"data":{}}`,
			want: "{\n\t\"data\": {}\n}",
		},
		{
			json: `{"array":[null,null]}`,
			want: "{\n\t\"array\": [\n\t\tnull,\n\t\tnull\n\t]\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			json := jsonFormatter(t, tt.json)
			t.Logf("\n%s", json)
			assert.Equal(t, tt.want, json)
		})
	}

	t.Run("error", func(t *testing.T) {
		tb := &bufferTB{name: t.Name()}
		assert.Panics(t, func() {
			jsonFormatter(tb, "")
		})
		_goldie.SetTest(t).Equal(tb.Bytes()).FailNow()
	})
}

func TestTool_JSONEq(t *testing.T) {
	tests := []struct {
		name   string
		got    string
		want   string
		failed bool
	}{
		{
			name:   "Succeeded",
			got:    "{}",
			want:   `{}`,
			failed: false,
		},
		{
			name:   "Failed",
			got:    "{}",
			want:   `{"data":null}`,
			failed: true,
		},
		{
			name:   "unexpected end of JSON input",
			got:    "",
			failed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &bufferTB{name: t.Name()}
			tl := SetTest(tb)
			tl.readFile = helperOSReadFile(t, []byte(tt.want), nil)

			cl := tl.JSONEq(tt.got)
			cl.Fail()
			if cl.Failed() {
				assert.Panics(t, func() { cl.FailNow() })
			} else {
				assert.NotPanics(t, func() { cl.FailNow() })
			}
			assert.Equal(t, tt.failed, cl.Failed())
			_goldie.SetTest(t).Equal(tb.Bytes()).Fail()
		})
	}

	t.Run("check-for-updates", func(t *testing.T) {
		got := []byte("{}")
		tb := &bufferTB{name: t.Name()}
		tl := SetTest(tb)
		tl.flag = new(bool)
		*tl.flag = true
		tl.readFile = helperOSReadFile(t, got, nil)
		tl.writeFile = func(name string, data []byte, mode os.FileMode) error {
			assert.Equal(t, name, "testdata/TestTool_JSONEq/check-for-updates.json.golden")
			assert.Equal(t, got, data)
			assert.Equal(t, tl.fileMode, mode)
			return nil
		}
		tl.mkdirAll = func(string, os.FileMode) error { return nil }

		cl := tl.JSONEq(string(got))
		cl.Fail()
		if cl.Failed() {
			assert.Panics(t, func() { cl.FailNow() })
		} else {
			assert.NotPanics(t, func() { cl.FailNow() })
		}
		assert.Equal(t, false, cl.Failed())
		_goldie.SetTest(t).Equal(tb.Bytes()).Fail()
	})
}

func TestJSONEq(t *testing.T) {
	tests := []struct {
		name   string
		got    string
		want   string
		failed bool
	}{
		{
			name:   "Succeeded",
			got:    "{}",
			want:   `{}`,
			failed: false,
		},
		{
			name:   "Failed",
			got:    "{}",
			want:   `{"data":null}`,
			failed: true,
		},
		{
			name:   "unexpected end of JSON input",
			got:    "",
			failed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin := _golden
			defer func() { _golden = origin }()

			tb := &bufferTB{name: t.Name()}
			_golden.readFile = helperOSReadFile(t, []byte(tt.want), nil)

			cl := JSONEq(tb, tt.got)
			cl.Fail()
			if cl.Failed() {
				assert.Panics(t, func() { cl.FailNow() })
			} else {
				assert.NotPanics(t, func() { cl.FailNow() })
			}
			assert.Equal(t, tt.failed, cl.Failed())
			_goldie.SetTest(t).Equal(tb.Bytes()).Fail()
		})
	}
}
