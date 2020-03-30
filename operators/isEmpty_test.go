package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestIsEmpty(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B").Pipe(operators.IsEmpty()),
			rx.Just("A").Pipe(operators.IsEmpty()),
			rx.Empty().Pipe(operators.IsEmpty()),
			rx.Throw(ErrTest).Pipe(operators.IsEmpty()),
		},
		[][]interface{}{
			{false, rx.Complete},
			{false, rx.Complete},
			{true, rx.Complete},
			{ErrTest},
		},
	)
}
