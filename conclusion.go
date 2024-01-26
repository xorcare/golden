// Copyright (c) 2019-2023 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package golden

import "fmt"

// Conclusion interface wrapping conclusion.
type Conclusion interface {
	// Failed reports whether the function has failed.
	Failed() bool
	// Fail marks the function as having failed but continues execution.
	// Also accompanying messages will be printed in the output of the test.
	// ATTENTION! executed only if expression is false `Failed() == true`.
	Fail()
	// FailNow marks the function as having failed and stops its execution
	// by calling runtime.Goexit (which then runs all deferred calls in the
	// current goroutine).
	// ATTENTION! executed only if expression is false `Failed() == true`.
	FailNow()
}

type conclusion struct {
	successful bool
	t          TestingTB
	diff       fmt.Stringer
}

func newConclusion(test TestingTB) conclusion {
	return conclusion{t: test}
}

// Failed reports whether the function has failed.
func (c conclusion) Failed() bool {
	return !c.successful
}

// Fail marks the function as having failed but continues execution.
// Also accompanying messages will be printed in the output of the test.
// ATTENTION! executed only if expression is false `Failed() == true`.
func (c conclusion) Fail() {
	if c.Failed() {
		c.t.Logf("%s", c.diff)
		c.t.Fail()
	}
}

// FailNow marks the function as having failed and stops its execution
// by calling runtime.Goexit (which then runs all deferred calls in the
// current goroutine).
// ATTENTION! executed only if expression is false `Failed() == true`.
func (c conclusion) FailNow() {
	if c.Failed() {
		c.t.Logf("%s", c.diff)
		c.t.FailNow()
	}
}
