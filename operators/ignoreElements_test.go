package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestIgnoreElements(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(operators.IgnoreElements()),
			rx.Just("A", "B", "C").Pipe(operators.IgnoreElements()),
			rx.Concat(rx.Just("A", "B", "C"), rx.Throw(ErrTest)).Pipe(operators.IgnoreElements()),
		},
		[][]interface{}{
			{rx.Completed},
			{rx.Completed},
			{ErrTest},
		},
	)
}
