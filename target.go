package golden

import "fmt"

var _ fmt.Stringer = target(0)

const (
	// Golden file target.
	Golden target = iota
	// Input file target.
	Input
	// latest the maximum target used. Should not be used in your code.
	latest
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
