package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestCatch(t *testing.T) {
	op := operators.Catch(
		func(error) rx.Observable {
			return rx.Just("D", "E")
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C").Pipe(op),
			rx.Concat(rx.Just("A", "B", "C"), rx.Throw(ErrTest)).Pipe(op),
		},
		[][]interface{}{
			{"A", "B", "C", rx.Completed},
			{"A", "B", "C", "D", "E", rx.Completed},
		},
	)
}
