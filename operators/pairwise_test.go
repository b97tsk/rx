package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestPairwise(t *testing.T) {
	op := rx.Pipe(operators.Pairwise(), ToString())
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(op),
			rx.Just("A").Pipe(op),
			rx.Just("A", "B").Pipe(op),
			rx.Just("A", "B", "C").Pipe(op),
			rx.Just("A", "B", "C", "D").Pipe(op),
			rx.Concat(rx.Just("A", "B", "C", "D"), rx.Throw(ErrTest)).Pipe(op),
		},
		[][]interface{}{
			{rx.Complete},
			{rx.Complete},
			{"[A B]", rx.Complete},
			{"[A B]", "[B C]", rx.Complete},
			{"[A B]", "[B C]", "[C D]", rx.Complete},
			{"[A B]", "[B C]", "[C D]", ErrTest},
		},
	)
}
