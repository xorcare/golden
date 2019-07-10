// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golden

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"unicode/utf8"
)

type target = uint

const (
	// Golden file target.
	Golden target = iota
	// Input file target.
	Input
)

// TestingTB is the interface common to T and B.
type TestingTB interface {
	Name() string
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Tool implements the basic logic of working with golden files.
// All functionality is implemented through a non-mutating state
// machine, which at a certain point in time can perform an action
// on the basis of the state or change any parameter by creating
// a new copy of the state.
type Tool struct {
	test     TestingTB
	dir      string
	fileMode os.FileMode
	inpExt   string
	modeDir  os.FileMode
	outExt   string
	target   target
	flag     *bool
	prefix   string

	mkdirAll  func(path string, perm os.FileMode) error
	readFile  func(filename string) ([]byte, error)
	remove    func(name string) error
	stat      func(name string) (os.FileInfo, error)
	writeFile func(filename string, data []byte, perm os.FileMode) error
}

// tool object with predefined parameters intended for use in global
// functions that provide a simplified api for convenient interaction
// with the functionality of the package.
var tool = Tool{
	// dir testdata is the directory for test data already accepted
	// in the standard library which is also ignored by standard
	// go tools and should not change in your tests.
	dir:      "testdata",
	outExt:   "golden",
	inpExt:   "input",
	fileMode: 0644,
	modeDir:  0755,
	target:   Golden,

	mkdirAll:  os.MkdirAll,
	readFile:  ioutil.ReadFile,
	remove:    os.Remove,
	stat:      os.Stat,
	writeFile: ioutil.WriteFile,
}

func init() {
	tool.flag = flag.Bool("update", false, "update test golden files")
}

// Assert is a tool to compare the actual value obtained in the test and
// the value from the golden file. Also built-in functionality for
// updating golden files using the command line flag.
func Assert(t TestingTB, got []byte) {
	tool.SetTest(t).Assert(got)
}

// Read is a functional for reading both input and golden files using
// the appropriate target.
func Read(t TestingTB) []byte {
	return tool.SetTest(t).SetTarget(Input).Read()
}

// Run is a functional that automates the process of reading the input file
// of the test bytes and the execution of the input function of testing and
// checking the results.
func Run(t TestingTB, do func(input []byte) (outcome []byte, err error)) {
	tool.SetTest(t).Run(do)
}

// SetTest is a mechanism to create a new copy of the base Tool object for
// advanced use. This method replaces the constructor for the Tool structure.
func SetTest(t TestingTB) Tool {
	return tool.SetTest(t)
}

// Assert is a tool to compare the actual value obtained in the test and
// the value from the golden file. Also built-in functionality for
// updating golden files using the command line flag.
func (tool Tool) Assert(got []byte) {
	tool.Update(got)
	tool.compare(got, tool.SetTarget(Golden).Read())
}

// Read is a functional for reading both input and golden files using
// the appropriate target.
func (tool Tool) Read() (bs []byte) {
	const f = "golden: read the value of nil since it is not found file: %s"

	bs, err := tool.readFile(tool.path())
	if os.IsNotExist(err) {
		tool.test.Logf(f, tool.path())
		return nil
	} else if err != nil {
		tool.test.Fatalf("golden: %s", err)
	}

	return bs
}

// Run is a functional that automates the process of reading the input file
// of the test bytes and the execution of the input function of testing and
// checking the results.
func (tool Tool) Run(do func(input []byte) (got []byte, err error)) {
	bs, err := do(tool.SetTarget(Input).Read())
	tool.ok(err)
	tool.Assert(bs)
}

// SetPrefix a prefix value setter.
func (tool Tool) SetPrefix(prefix string) Tool {
	tool.prefix = prefix
	return tool
}

// SetTarget a target value setter.
func (tool Tool) SetTarget(tar target) Tool {
	tool.target = tar
	return tool
}

// SetTest a test value setter in the call chain must be used first
// to prevent abnormal situations when using other methods.
func (tool Tool) SetTest(t TestingTB) Tool {
	tool.test = t
	return tool
}

// Update functional reviewer is the need to update the golden files
// and doing it.
func (tool Tool) Update(bs []byte) {
	if tool.flag == nil || !*tool.flag {
		return
	}

	tool.test.Logf("golden: updating file: %s", tool.path())
	tool.write(bs)
}

// write is a functional for writing both input and golden files using
// the appropriate target.
func (tool Tool) write(bs []byte) {
	path := tool.path()
	tool.mkdir(filepath.Dir(path))
	tool.test.Logf("golden: start write to file: %s", path)
	if bs == nil {
		tool.test.Logf("golden: nil value will not be written")
		fileInfo, err := tool.stat(path)
		if err == nil && !fileInfo.IsDir() {
			tool.test.Logf("golden: current test bytes file will be deleted")
			tool.ok(tool.remove(path))
		}
		if !os.IsNotExist(err) {
			tool.ok(err)
		}
	} else {
		tool.ok(tool.writeFile(path, bs, tool.fileMode))
	}
}

// compare mechanism to compare the bytes.
func (tool Tool) compare(got, want []byte) {
	if !bytes.Equal(got, want) {
		format := "golden: compare error got = %#v, want %#v"
		if utf8.ValidString(string(want)) || utf8.ValidString(string(got)) {
			format = "golden: compare error got = %q, want %q"
		}

		tool.test.Fatalf(format, got, want)
	}
}

// mkdir the mechanism to create the directory.
func (tool Tool) mkdir(loc string) {
	fileInfo, err := tool.stat(loc)
	switch {
	case err != nil && os.IsNotExist(err):
		err = tool.mkdirAll(loc, tool.modeDir)
	case err == nil && !fileInfo.IsDir():
		tool.test.Errorf("golden: test dir is a file: %s", loc)
	}

	tool.ok(err)
}

// ok fails the test if an err is not nil.
func (tool Tool) ok(err error) {
	if err != nil {
		tool.test.Fatalf("golden: %s", err)
	}
}

// path is getter to get the path to the file containing the test data.
func (tool Tool) path() (path string) {
	ext := tool.outExt
	if tool.target == Input {
		ext = tool.inpExt
	}

	s := fmt.Sprintf("%s.%s", tool.test.Name(), ext)
	if tool.prefix != "" {
		s = fmt.Sprintf("%s.%s.%s", tool.test.Name(), tool.prefix, ext)
	}

	return filepath.Join(tool.dir, s)
}
