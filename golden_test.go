// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

var helper = tool

func TestMain(m *testing.M) {
	helper.UpdateFlag = tool.UpdateFlag
	tool.UpdateFlag = nil
	os.Exit(m.Run())
}

func TestToolDefault(t *testing.T) {
	jsonBytes, err := json.MarshalIndent(tool, "", "\t")
	if err != nil {
		t.Fatalf("TestToolDefault() failed json.Marshal(%#v), \n error: %v", tool, err)
	}
	helper.SetTest(t).Assert(jsonBytes)
}

func TestAssert(t *testing.T) {
	type args struct {
		test *FakeTest
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
				test: new(FakeTest),
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
				test: new(FakeTest),
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
				test: new(FakeTest),
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
				test: new(FakeTest),
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
			origin := tool
			defer func() { tool = origin }()

			tool.readFile = func(filename string) (bytes []byte, e error) {
				t.Logf(`os.ReadFile(%q) `, filename)
				return tt.readFile.bytes, tt.readFile.error
			}
			{
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					}
					tt.args.test.Assert(t)
				}()
				tt.args.test.name = t.Name()
				Assert(tt.args.test, tt.args.got)
			}
		})
	}
}

func TestRead(t *testing.T) {
	type args struct {
		test *FakeTest
		tar  target
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
				test: new(FakeTest),
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
				test: new(FakeTest),
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
				test: new(FakeTest),
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
			origin := tool
			defer func() { tool = origin }()

			tool.readFile = func(filename string) (bytes []byte, e error) {
				t.Logf(`os.ReadFile(%q) `, filename)
				return tt.readFile.bytes, tt.readFile.error
			}
			{
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					}
					tt.args.test.Assert(t)
				}()
				tt.args.test.name = t.Name()
				got := Read(tt.args.test, tt.args.tar)
				if !bytes.Equal(got, tt.want) {
					t.Errorf("Read() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		test *FakeTest
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
				test: new(FakeTest),
				do: func(input []byte) (outcome []byte, err error) {
					return nil, nil
				},
			},
			recover: false,
		},
		{
			name: "run-with-error",
			args: args{
				test: new(FakeTest),
				do: func(input []byte) (outcome []byte, err error) {
					return nil, os.ErrClosed
				},
			},
			recover: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin := tool
			defer func() { tool = origin }()

			tool.readFile = func(filename string) (bytes []byte, e error) {
				t.Logf(`os.ReadFile(%q)`, filename)
				return nil, nil
			}
			{
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					}
					tt.args.test.Assert(t)
				}()
				tt.args.test.name = t.Name()
				Run(tt.args.test, tt.args.do)
			}
		})
	}
}

func TestSetTest(t *testing.T) {
	type args struct {
		test tb
	}
	m := new(FakeTest)
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
			want: Tool{Test: m},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origin := tool
			defer func() { tool = origin }()

			if got := SetTest(tt.args.test); got.Test != tt.want.Test {
				t.Errorf("SetTest() = %v, want %v", got.Test, tt.want.Test)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	type args struct {
		test  *FakeTest
		tar   target
		bytes []byte
	}
	type stat struct {
		fileInfo *FakeFile
		error    error
	}
	tests := []struct {
		name      string
		args      args
		writeFile error
		stat      stat
		recover   bool
	}{
		{
			name: "write-nil",
			args: args{
				test:  new(FakeTest),
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
				test:  new(FakeTest),
				tar:   Golden,
				bytes: nil,
			},
			stat: stat{
				fileInfo: new(FakeFile),
			},
			recover: false,
		},
		{
			name: "write-empty",
			args: args{
				test:  new(FakeTest),
				tar:   Golden,
				bytes: []byte{},
			},
			stat: stat{
				fileInfo: new(FakeFile),
			},
			recover: false,
		},
		{
			name: "write-bytes",
			args: args{
				test:  new(FakeTest),
				tar:   Golden,
				bytes: []byte("golden"),
			},
			stat: stat{
				fileInfo: new(FakeFile),
			},
			recover: false,
		},
		{
			name: "fatality-error",
			args: args{
				test:  new(FakeTest),
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
			origin := tool
			defer func() { tool = origin }()

			tool.writeFile = func(filename string, data []byte, perm os.FileMode) error {
				t.Logf(`os.WriteFile(%q, %q, %d) `, filename, data, perm)
				return tt.writeFile
			}
			tool.mkdirAll = func(path string, perm os.FileMode) error {
				t.Logf(`os.MkdirAll(%q, %d) `, path, perm)
				return nil
			}
			tool.remove = func(name string) error {
				t.Logf(`os.Remove(%q)`, name)
				return nil
			}
			tool.stat = func(name string) (os.FileInfo, error) {
				t.Logf(`os.Stat(%q)`, name)
				if tt.stat.fileInfo != nil {
					tt.stat.fileInfo.name = name
				}
				return tt.stat.fileInfo, tt.stat.error
			}
			{
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					}
					tt.args.test.Assert(t)
				}()
				tt.args.test.name = t.Name()
				Write(tt.args.test, tt.args.tar, tt.args.bytes)
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
		test     FakeTest
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
			{
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					}
					tt.test.Assert(t)
				}()
				tt.test.name = t.Name()
				tt.tool.SetTest(&tt.test).Assert(tt.args.got)
			}
		})
	}
}

func TestTool_Path(t *testing.T) {
	tests := []struct {
		name     string
		tool     Tool
		wantPath string
	}{
		{
			name:     "initial",
			tool:     Tool{},
			wantPath: "TestTool_Path/initial.",
		},
		{
			name:     "input-target-path",
			tool:     tool.SetTarget(Input),
			wantPath: "testdata/TestTool_Path/input-target-path.input",
		},
		{
			name:     "golden-target-path",
			tool:     tool.SetTarget(Golden),
			wantPath: "testdata/TestTool_Path/golden-target-path.golden",
		},
		{
			name:     "default-target-path",
			tool:     tool.SetTarget(0),
			wantPath: "testdata/TestTool_Path/default-target-path.input",
		},
		{
			name:     "set-index-path",
			tool:     tool.SetIndex(1).SetTarget(Golden),
			wantPath: "testdata/TestTool_Path/set-index-path.001.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tool.Test = &FakeTest{name: t.Name()}
			if gotPath := tt.tool.Path(); gotPath != tt.wantPath {
				t.Errorf("Tool.Path() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func TestTool_Read(t *testing.T) {
	type args struct {
		test *FakeTest
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
				test: new(FakeTest),
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
				test: new(FakeTest),
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
				test: new(FakeTest),
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
			{
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					}
					tt.args.test.Assert(t)
				}()
				got := tt.tool.SetTest(tt.args.test).SetTarget(tt.args.tar).Read()
				if !bytes.Equal(got, tt.want) {
					t.Errorf("Read() = %v, want %v", got, tt.want)
				}
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
		test    FakeTest
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
			{
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					}
					tt.test.Assert(t)
				}()
				tt.tool.SetTest(&tt.test).Run(tt.args.do)
			}
		})
	}
}

func TestTool_SetDir(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		tool Tool
		args args
		want Tool
	}{
		{
			name: "golden",
			args: args{
				dir: "golden",
			},
			want: Tool{
				Dir: "golden",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tool.SetDir(tt.args.dir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tool.SetDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTool_SetIndex(t *testing.T) {
	type args struct {
		index uint8
	}
	tests := []struct {
		name string
		tool Tool
		args args
		want Tool
	}{
		{
			name: "set-one",
			args: args{
				index: 1,
			},
			want: Tool{
				Index: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tool.SetIndex(tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tool.SetIndex() = %v, want %v", got, tt.want)
			}
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
				Target: Input,
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
		t tb
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
			want: Tool{Test: t},
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
		fileInfo *FakeFile
		error    error
	}
	fal := false
	tru := true
	tests := []struct {
		name string
		tool Tool
		test FakeTest
		args args
		stat stat
	}{
		{
			name: "not-update-with-nil",
			tool: Tool{UpdateFlag: nil},
			args: args{
				[]byte("golden"),
			},
		},
		{
			name: "not-update-with-false",
			tool: Tool{UpdateFlag: &fal},
			args: args{
				[]byte("golden"),
			},
		},
		{
			name: "update-with-true",
			tool: Tool{UpdateFlag: &tru},
			args: args{
				[]byte("golden"),
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

func TestTool_Write(t *testing.T) {
	type args struct {
		test  *FakeTest
		tar   target
		bytes []byte
	}
	type stat struct {
		fileInfo *FakeFile
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
				test:  new(FakeTest),
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
				test:  new(FakeTest),
				tar:   Golden,
				bytes: nil,
			},
			stat: stat{
				fileInfo: new(FakeFile),
				error:    nil,
			},
			recover: false,
		},
		{
			name: "write-empty",
			args: args{
				test:  new(FakeTest),
				tar:   Golden,
				bytes: []byte{},
			},
			stat: stat{
				fileInfo: new(FakeFile),
			},
			recover: false,
		},
		{
			name: "write-bytes",
			args: args{
				test:  new(FakeTest),
				tar:   Golden,
				bytes: []byte("golden"),
			},
			stat: stat{
				fileInfo: new(FakeFile),
			},
			recover: false,
		},
		{
			name: "fatality-error",
			args: args{
				test:  new(FakeTest),
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
			{
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					}
					tt.args.test.Assert(t)
				}()
				tt.tool.SetTest(tt.args.test).
					SetTarget(tt.args.tar).
					Write(tt.args.bytes)
			}
		})
	}
}

func TestTool_compare(t *testing.T) {
	type args struct {
		got  []byte
		want []byte
	}
	tests := []struct {
		name    string
		tool    Tool
		test    FakeTest
		args    args
		recover bool
	}{
		{
			name: "equal-nil",
			args: args{
				got:  nil,
				want: nil,
			},
			recover: false,
		},
		{
			name: "equal-bytes",
			args: args{
				got:  []byte("golden"),
				want: []byte("golden"),
			},
			recover: false,
		},
		{
			name: "not-equal-nil-and-bytes",
			args: args{
				got:  []byte("golden"),
				want: nil,
			},
			recover: true,
		},
		{
			name: "not-equal-bytes",
			args: args{
				got:  []byte("golden"),
				want: []byte("Z29sZGVu"),
			},
			recover: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			{
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					}
					tt.test.Assert(t)
				}()
				tt.test.name = t.Name()
				tt.tool.SetTest(&tt.test).compare(tt.args.got, tt.args.want)
			}
		})
	}
}

func TestTool_mkdir(t *testing.T) {
	type args struct {
		loc string
	}
	type stat struct {
		fileInfo *FakeFile
		error    error
	}
	tests := []struct {
		name     string
		tool     Tool
		test     FakeTest
		args     args
		stat     stat
		mkdirAll error
		recover  bool
	}{
		{
			name: "fatality-error",
			args: args{
				loc: tool.SetTest(t).Path(),
			},
			stat: stat{
				error: os.ErrPermission,
			},
			recover: true,
		},
		{
			name: "error-file-does-not-exist",
			args: args{
				loc: tool.SetTest(t).Path(),
			},
			stat: stat{
				error: os.ErrNotExist,
			},
			recover: false,
		},
		{
			name: "error-dir-is-a-file",
			args: args{
				loc: tool.SetTest(t).Path(),
			},
			stat: stat{
				fileInfo: new(FakeFile),
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
			{
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					}
					tt.test.Assert(t)
				}()
				tt.test.name = t.Name()
				tt.tool.SetTest(&tt.test).mkdir(tt.args.loc)
			}
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
		test    FakeTest
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
			{
				defer func() {
					if r := recover(); (r == nil) == tt.recover {
						t.Error(r)
					}
					tt.test.Assert(t)
				}()
				tt.test.name = t.Name()
				tt.tool.SetTest(&tt.test).ok(tt.args.err)
			}
		})
	}
}

// FakeTest implements tb methods.
type FakeTest struct {
	name string
	Errs []string `json:"errs,omitempty"`
	Logs []string `json:"logs,omitempty"`
	Fats []string `json:"fats,omitempty"`
}

// tb interface methods.

func (m FakeTest) Name() string {
	return m.name
}

func (m *FakeTest) Logf(format string, args ...interface{}) {
	m.Logs = append(m.Logs, fmt.Sprintf(format, args...))
}

func (m *FakeTest) Errorf(format string, args ...interface{}) {
	m.Errs = append(m.Errs, fmt.Sprintf(format, args...))
}

func (m *FakeTest) Fatalf(format string, args ...interface{}) {
	m.Fats = append(m.Fats, fmt.Sprintf(format, args...))
	panic(fmt.Sprintf(format, args...))
}

// test control methods.

func (m *FakeTest) Assert(t tb) {
	jsonBytes, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		t.Fatalf("FakeTest.Assert() failed json.Marshal(%#v), error: %v", m, err)
	}
	helper.SetTest(t).Assert(jsonBytes)
}

// FakeFile implements FileLike and also os.FileInfo.
type FakeFile struct {
	name     string
	contents string
	mode     os.FileMode
	offset   int
}

// FileLike methods.

func (f *FakeFile) Name() string {
	// A bit of a cheat: we only have a basename, so that's also ok for FileInfo.
	return f.name
}

// os.FileInfo methods.

func (f *FakeFile) Size() int64 {
	return int64(len(f.contents))
}

func (f *FakeFile) Mode() os.FileMode {
	return f.mode
}

func (f *FakeFile) ModTime() time.Time {
	return time.Time{}
}

func (f *FakeFile) IsDir() bool {
	return false
}

func (f *FakeFile) Sys() interface{} {
	return nil
}
