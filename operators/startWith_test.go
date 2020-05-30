package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestStartWith(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("D", "E").Pipe(operators.StartWith("A", "B", "C")),
			rx.Empty().Pipe(operators.StartWith("A", "B", "C")),
			rx.Throw(ErrTest).Pipe(operators.StartWith("A", "B", "C")),
			rx.Throw(ErrTest).Pipe(operators.StartWith()),
		},
		[][]interface{}{
			{"A", "B", "C", "D", "E", rx.Completed},
			{"A", "B", "C", rx.Completed},
			{"A", "B", "C", ErrTest},
			{ErrTest},
		},
	)
}
