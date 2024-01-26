// Copyright (c) 2019-2023 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golden_test

import (
	"encoding/base64"
	"fmt"
	"path"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/stretchr/testify/assert"

	"github.com/xorcare/golden"
)

// In the current example, the value of `got` will be compared with the value
// of `want` that we get from the file `testdata/ExampleAssert.golden` and
// after comparing the data if there is a difference, the test will be aborted
// with an error. If you run the test flag is used with the `-update` data from
// the variable `got` is written in the golden file.
//
// The test name is assumed to be equal to ExampleAssert.
func ExampleAssert() {
	t := newTestingT()

	got, err := base64.RawURLEncoding.DecodeString("Z29sZGVu")
	assert.NoError(t, err)

	golden.Assert(t, got)

	// Output:
	// golden: read the value of nil since it is not found file: testdata/TestExamples/ExampleAssert.golden
	//
	// Error Trace:
	// Error:      Not equal:
	//             expected: "[]byte(nil)"
	//             actual  : "golden"
	//
	//             Diff:
	//             --- Expected
	//             +++ Actual
	//             @@ -1 +1 @@
	//             -[]byte(nil)
	//             +golden

}

// In the current example, when you run the Run function, the data from the
// `testdata/ExampleRun.input` file will be read and with this data will be
// called the function, which is passed in, as a result of the function
// execution we will get the `got` data and possibly an error, which will
// be processed by the internal method implementation. The `got` will be
// compared with the meaning `want` which we get from the file
// `test data/ExampleRun.golden` and after data comparison in case of
// differences, the test will be fail. If you run the test flag is used with
// the `-update` data from the variable `got` is written in the golden file.
//
// The test name is assumed to be equal to ExampleRun.
func ExampleRun() {
	t := newTestingT()

	golden.Run(t, func(input []byte) (got []byte, err error) {
		return base64.RawURLEncoding.DecodeString(string(input))
	})
	// Output:
	// golden: read the value of nil since it is not found file: testdata/TestExamples/ExampleRun.input
	// golden: read the value of nil since it is not found file: testdata/TestExamples/ExampleRun.golden
	//
	// Error Trace:
	// Error:      Not equal:
	//             expected: "[]byte(nil)"
	//             actual  : ""
	//
	//             Diff:
	//             --- Expected
	//             +++ Actual
	//             @@ -1 +1 @@
	//             -[]byte(nil)
	//             +
}

// ExampleRead the example shows how you can use the global api to read files
// together with the already considered golden.Assert.
//
// The test name is assumed to be equal to ExampleRead.
func ExampleRead() {
	t := newTestingT()

	input := string(golden.Read(t))
	got, err := base64.RawURLEncoding.DecodeString(input)
	assert.NoError(t, err)

	golden.Assert(t, got)

	// Output:
	// golden: read the value of nil since it is not found file: testdata/TestExamples/ExampleRead.input
	// golden: read the value of nil since it is not found file: testdata/TestExamples/ExampleRead.golden
	//
	// Error Trace:
	// Error:      Not equal:
	//             expected: "[]byte(nil)"
	//             actual  : ""
	//
	//             Diff:
	//             --- Expected
	//             +++ Actual
	//             @@ -1 +1 @@
	//             -[]byte(nil)
	//             +
}

type T struct {
	name string
}

func (t *T) Fail()    {}
func (t *T) FailNow() {}
func (t *T) Helper()  {}

func (t T) Name() string { return t.name }

func (t *T) Errorf(format string, args ...interface{}) { t.Logf(format, args...) }
func (t *T) Fatalf(format string, args ...interface{}) { t.Logf(format, args...) }

func (t *T) Logf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	// Removed the trace, it contains line numbers and changes dynamically,
	// it is not convenient to see in the examples.
	re := regexp.MustCompile(`(?im)^\t?Error\ Trace\:([\S\s\n]+)^\t?Error\:`)
	msg = re.ReplaceAllString(msg, "\tError Trace:\n\tError:")

	msg = strings.Replace(msg, "\t", "", -1)

	// Trimming lines consisting only of spaces or containing spaces to the right.
	re = regexp.MustCompile(`(?im)^(.*)$`)
	msg = re.ReplaceAllStringFunc(msg, func(s string) string {
		return strings.TrimRightFunc(s, unicode.IsSpace)
	})

	fmt.Println(msg)
}

func newTestingT() *T {
	t := T{name: "TestExamples"}
	t.name = path.Join(t.name, caller(2))
	return &t
}

func caller(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		panic(fmt.Sprintf("Couldn't get the caller info level: %d", skip))
	}

	fp := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fp, ".")
	fn := parts[len(parts)-1]

	return fn
}
