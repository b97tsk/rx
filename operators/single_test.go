package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestSingle(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B").Pipe(operators.Single()),
			rx.Just("A").Pipe(operators.Single()),
			rx.Empty().Pipe(operators.Single()),
			rx.Throw(ErrTest).Pipe(operators.Single()),
		},
		[][]interface{}{
			{rx.ErrNotSingle},
			{"A", rx.Completed},
			{rx.ErrEmpty},
			{ErrTest},
		},
	)
}
