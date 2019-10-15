// Copyright © 2019, Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golden

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/stretchr/testify/assert"
)

type target uint

func (t target) String() string {
	switch t {
	case Golden:
		return "golden"
	case Input:
		return "input"
	default:
		panic(fmt.Sprintf("unsupported target: %d", t))
	}
}

const (
	// Golden file target.
	Golden target = iota
	// Input file target.
	Input
	// latest the maximum target used. Should not be used in your code.
	latest
)

// TestingTB is the interface common to T and B.
type TestingTB interface {
	Name() string
	Logf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	FailNow()
	Fail()
}

type testingHelper interface {
	Helper()
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
	modeDir  os.FileMode
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
var _golden = Tool{
	// dir testdata is the directory for test data already accepted
	// in the standard library, which is also ignored by standard
	// go tools and should not change in your tests.
	dir:      "testdata",
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
	_golden.flag = flag.Bool("update", false, "update test golden files")
}

// Assert is a tool to compare the actual value obtained in the test and
// the value from the golden file. Also, built-in functionality for
// updating golden files using the command line flag.
func Assert(t TestingTB, got []byte) {
	if h, ok := t.(testingHelper); ok {
		h.Helper()
	}
	SetTest(t).Assert(got)
}

// Equal is a tool to compare the actual value obtained in the test and
// the value from the golden file. Also, built-in functionality for
// updating golden files using the command line flag.
func Equal(t TestingTB, got []byte) Conclusion {
	if h, ok := t.(testingHelper); ok {
		h.Helper()
	}
	return SetTest(t).Equal(got)
}

// Read is a functional for reading both input and golden files using
// the appropriate target.
func Read(t TestingTB) []byte {
	return SetTest(t).SetTarget(Input).Read()
}

// Run is a functional that automates the process of reading the input file
// of the test bytes and the execution of the input function of testing and
// checking the results.
func Run(t TestingTB, do func(input []byte) (outcome []byte, err error)) {
	SetTest(t).Run(do)
}

// SetTest is a mechanism to create a new copy of the base Tool object for
// advanced use. This method replaces the constructor for the Tool structure.
func SetTest(t TestingTB) Tool {
	return _golden.SetTest(t)
}

// Assert is a tool to compare the actual value obtained in the test and
// the value from the golden file. Also, built-in functionality for
// updating golden files using the command line flag.
func (t Tool) Assert(got []byte) {
	t.Update(got)
	if h, ok := t.test.(testingHelper); ok {
		h.Helper()
	}
	t.Equal(got).FailNow()
}

// Equal is a tool to compare the actual value obtained in the test and
// the value from the golden file. Also, built-in functionality for
// updating golden files using the command line flag.
func (t Tool) Equal(actual []byte) Conclusion {
	t.Update(actual)
	if h, ok := t.test.(testingHelper); ok {
		h.Helper()
	}

	expected := t.SetTarget(Golden).Read()

	if expected == nil {
		expected = []byte(fmt.Sprintf("%#v", expected))
	}
	if actual == nil {
		actual = []byte(fmt.Sprintf("%#v", actual))
	}

	i := new(interceptor)
	c := newConclusion(t.test)
	c.successful = assert.Equal(i, string(expected), string(actual))
	c.diff = i

	return c
}

// Read is a functional for reading both input and golden files using
// the appropriate target.
func (t Tool) Read() (bs []byte) {
	const f = "golden: read the value of nil since it is not found file: %s"

	bs, err := t.readFile(t.path())
	if os.IsNotExist(err) {
		t.test.Logf(f, t.path())
		return nil
	} else if err != nil {
		t.test.Fatalf("golden: %s", err)
	}

	return bs
}

// Run is a functional that automates the process of reading the input file
// of the test bytes and the execution of the input function of testing and
// checking the results.
func (t Tool) Run(do func(input []byte) (got []byte, err error)) {
	bs, err := do(t.SetTarget(Input).Read())
	t.ok(err)
	t.Assert(bs)
}

// SetPrefix a prefix value setter.
func (t Tool) SetPrefix(prefix string) Tool {
	t.prefix = rewrite(prefix)
	return t
}

// SetTarget a target value setter.
func (t Tool) SetTarget(tar target) Tool {
	t.target = tar
	return t
}

// SetTest a test value setter in the call chain must be used first
// to prevent abnormal situations when using other methods.
func (t Tool) SetTest(tb TestingTB) Tool {
	t.test = tb
	return t
}

// Update functional reviewer is the need to update the golden files
// and doing it.
func (t Tool) Update(bs []byte) {
	if t.flag == nil || !*t.flag {
		return
	}

	t.test.Logf("golden: updating file: %s", t.path())
	t.write(bs)
}

// write is a functional for writing both input and golden files using
// the appropriate target.
func (t Tool) write(bs []byte) {
	path := t.path()
	t.mkdir(filepath.Dir(path))
	t.test.Logf("golden: start write to file: %s", path)
	if bs == nil {
		t.test.Logf("golden: nil value will not be written")
		fileInfo, err := t.stat(path)
		if err == nil && !fileInfo.IsDir() {
			t.test.Logf("golden: current test bytes file will be deleted")
			t.ok(t.remove(path))
		}
		if !os.IsNotExist(err) {
			t.ok(err)
		}
	} else {
		t.ok(t.writeFile(path, bs, t.fileMode))
	}
}

// mkdir the mechanism to create the directory.
func (t Tool) mkdir(loc string) {
	fileInfo, err := t.stat(loc)
	switch {
	case err != nil && os.IsNotExist(err):
		t.test.Logf("golden: trying to create a directory: %q", loc)
		err = t.mkdirAll(loc, t.modeDir)
	case err == nil && !fileInfo.IsDir():
		t.test.Errorf("golden: test dir is a file: %s", loc)
	}

	t.ok(err)
}

// ok fails the test if an err is not nil.
func (t Tool) ok(err error) {
	if err != nil {
		t.test.Fatalf("golden: %s", err)
	}
}

// path is getter to get the path to the file containing the test data.
func (t Tool) path() (path string) {
	format := "%s"
	args := []interface{}{t.test.Name()}

	if t.prefix != "" {
		args = append(args, t.prefix)
	}

	// Add a target extansion. Always added last.
	args = append(args, t.target.String())
	// We add placeholders for the number of parameters excluding the name
	// of the test to print all the parameters.
	format += strings.Repeat(".%s", len(args)-1)
	return filepath.Join(t.dir, fmt.Sprintf(format, args...))
}

// rewrite rewrites a subname to having only printable characters and no white
// space.
func rewrite(str string) string {
	bs := make([]byte, 0, len(str))
	for _, b := range str {
		switch {
		case unicode.IsSpace(b):
			bs = append(bs, '_')
		case !strconv.IsPrint(b):
			s := strconv.QuoteRune(b)
			bs = append(bs, s[1:len(s)-1]...)
		default:
			bs = append(bs, string(b)...)
		}
	}
	return string(bs)
}
