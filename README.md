# Golden

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

## Installation

```bash
go get github.com/xorcare/golden
```

## Examples

 * [golden.Assert](https://godoc.org/github.com/xorcare/golden#example-Assert)
 * [golden.Read](https://godoc.org/github.com/xorcare/golden#example-Read)
 * [golden.Run](https://godoc.org/github.com/xorcare/golden#example-Run)

## Inspiration

 * [Golden Files—Why you should use them](https://medium.com/@jarifibrahim/golden-files-why-you-should-use-them-47087ec994bf)
 * [Testing with golden files in Go](https://medium.com/soon-london/testing-with-golden-files-in-go-7fccc71c43d3)
 * [Go advanced testing tips & tricks](https://medium.com/@povilasve/go-advanced-tips-tricks-a872503ac859)

## License

© Vasiliy Vasilyuk, 2019-2021

Released under the [BSD 3-Clause License][LICENSE].

[LICENSE]: https://github.com/xorcare/golden/blob/master/LICENSE 'BSD 3-Clause "New" or "Revised" License'
