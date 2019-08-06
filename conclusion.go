package golden

//go:generate ifacemaker -f $GOFILE -p $GOPACKAGE -s conclusion -i Conclusion -o itf-conclusion.go -D -y "Conclusion interface wrapping conclusion."  -c "Code generated file. DO NOT EDIT."

type conclusion struct {
	successful bool
	t          TestingTB
	diff       interface{}
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
