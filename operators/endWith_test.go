package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestEndWith(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C").Pipe(operators.EndWith("D", "E")),
			rx.Empty().Pipe(operators.EndWith("D", "E")),
			rx.Throw(ErrTest).Pipe(operators.EndWith("D", "E")),
			rx.Throw(ErrTest).Pipe(operators.EndWith()),
		},
		[][]interface{}{
			{"A", "B", "C", "D", "E", Completed},
			{"D", "E", Completed},
			{ErrTest},
			{ErrTest},
		},
	)
}
