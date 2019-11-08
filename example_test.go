package golden_test

import (
	"encoding/base64"
	"testing"

	"github.com/xorcare/golden"
)

// It is necessary for the syntactic correctness of the examples.
var t = &testing.T{}

// In the current example, the value of `got` will be compared with the value
// of `want` that we get from the file `testdata/ExampleAssert.golden` and
// after comparing the data if there is a difference, the test will be aborted
// with an error. If you run the test flag is used with the `-update` data from
// the variable `got` is written in the golden file.
//
// The test name is assumed to be equal to ExampleAssert.
func ExampleAssert() {
	got, err := base64.RawURLEncoding.DecodeString("Z29sZGVu")
	if err != nil {
		t.Fatalf("%s", err)
	}

	golden.Assert(t, got)
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
	golden.Run(t, func(input []byte) (got []byte, err error) {
		return base64.RawURLEncoding.DecodeString(string(input))
	})
}

// ExampleRead the example shows how you can use the global api to read files
// together with the already considered golden.Assert.
//
// The test name is assumed to be equal to ExampleRead.
func ExampleRead() {
	input := string(golden.Read(t))
	got, err := base64.RawURLEncoding.DecodeString(input)
	if err != nil {
		t.Fatal(err)
	}

	golden.Assert(t, got)
}
