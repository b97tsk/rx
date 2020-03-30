package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Pairwise(t *testing.T) {
	op := Pipe(operators.Pairwise(), toString)
	subscribeN(
		t,
		[]Observable{
			Empty().Pipe(op),
			Just("A").Pipe(op),
			Just("A", "B").Pipe(op),
			Just("A", "B", "C").Pipe(op),
			Just("A", "B", "C", "D").Pipe(op),
			Concat(Just("A", "B", "C", "D"), Throw(errTest)).Pipe(op),
		},
		[][]interface{}{
			{Complete},
			{Complete},
			{"[A B]", Complete},
			{"[A B]", "[B C]", Complete},
			{"[A B]", "[B C]", "[C D]", Complete},
			{"[A B]", "[B C]", "[C D]", errTest},
		},
	)
}
