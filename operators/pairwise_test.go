package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
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
			{rx.Completed},
			{rx.Completed},
			{"[A B]", rx.Completed},
			{"[A B]", "[B C]", rx.Completed},
			{"[A B]", "[B C]", "[C D]", rx.Completed},
			{"[A B]", "[B C]", "[C D]", ErrTest},
		},
	)
}
