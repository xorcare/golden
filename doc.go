// Copyright (c) 2019-2024 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
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
variables of go this is very convenient in case of significant changes in the
behavior of the system but also requires attention to the changed test data and
checking the correctness of the new golden results.

A special cli is provided in the package. The special flag `-update` is
provided in the package for conveniently updating ethos files, for example,
using the following command:

	go test ./... -update

Golden files are placed in directory `testdata` this directory is ignored by
the standard tools go, and it can accommodate a variety of data used in test or
samples.
*/
package golden
