# Golden

[![Build Status](https://travis-ci.org/xorcare/golden.svg?branch=master)](https://travis-ci.org/xorcare/golden)
[![codecov](https://codecov.io/gh/xorcare/golden/badge.svg)](https://codecov.io/gh/xorcare/golden)
[![Go Report Card](https://goreportcard.com/badge/github.com/xorcare/golden)](https://goreportcard.com/report/github.com/xorcare/golden)
[![GoDoc](https://godoc.org/github.com/xorcare/golden?status.svg)](https://godoc.org/github.com/xorcare/golden)

Package golden testing with golden files in Go. A golden file is the expected
output of test, stored as a separate file rather than as a string literal inside
the test code. So when the test is executed, it will read data from the file and
compare it to the output produced by the functionality under test.

When writing unit tests, there comes a point when you need to check that the
complex output of a function matches your expectations. This could be binary
data (like an Image, JSON, HTML etc). Golden files are a way of ensuring your
function output matches its .golden file. It’s a pattern used in the Go standard
library.

One of the advantages of the gold files approach is that you can easily update
the test data using the command line flag without copying the data into the text
variables of go this is very convenient in case of  significant changes in the
behavior of the system but also requires attention to the changed test data and
checking the correctness of the new golden results.

A special cli is provided in the package. The special flag `-update` is
provided in the package for conveniently updating ethos files, for example,
using the following command: 

	go test ./... -update

Golden files are placed in directory `testdata` this directory is ignored by
the standard tools go, and it can accommodate a variety of data used in test or
samples.

## Examples

**Using** `golden.Assert` - to check and save the data in the golden file.

```go
func TestDecode(t *testing.T) {
    got, err := base64.RawURLEncoding.DecodeString("Z29sZGVu")
    if err != nil {
        t.Fatal(err)
    }

    golden.Assert(t, got)
}
```

**Using** `golden.Run` - to automatically read the input file, run the
function you created and compare the results with the golden files.

```go
func TestDecode(t *testing.T) {
    golden.Run(t, func(input []byte) (got []byte, err error) {
        return base64.RawURLEncoding.DecodeString(string(input))
    })
}
```

**Using** `golden.Read` - you can use golden code to create
your own codebase, such as reading from an input file.

```go
func TestDecode(t *testing.T) {
    input := string(golden.Read(t, golden.Input))
    got, err := base64.RawURLEncoding.DecodeString(input)
    if err != nil {
        t.Fatal(err)
    }
    
    golden.Assert(t, got)
}
```

## Inspiration

 * [Golden Files — Why you should use them](http://bit.ly/2JikzYp)
 * [Testing with golden files in Go](http://bit.ly/2TFjdvC)
 * [Go advanced testing tips & tricks](https://bit.ly/2Cpi28Q)

## License

© Vasiliy Vasilyuk, 2019

Released under the [BSD 3-Clause License][LICENSE].

[LICENSE]: https://git.io/fhjjx 'BSD 3-Clause "New" or "Revised" License'
