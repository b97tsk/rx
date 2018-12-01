package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Pairwise(t *testing.T) {
	op := Pipe(operators.Pairwise(), toString)
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(op),
			Just("A").Pipe(op),
			Just("A", "B").Pipe(op),
			Just("A", "B", "C").Pipe(op),
			Just("A", "B", "C", "D").Pipe(op),
			Concat(Just("A", "B", "C", "D"), Throw(xErrTest)).Pipe(op),
		},
		xComplete,
		xComplete,
		"[A B]", xComplete,
		"[A B]", "[B C]", xComplete,
		"[A B]", "[B C]", "[C D]", xComplete,
		"[A B]", "[B C]", "[C D]", xErrTest,
	)
}
