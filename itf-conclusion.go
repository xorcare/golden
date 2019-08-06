// Code generated file. DO NOT EDIT.

package golden

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
