package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestFindIndex(t *testing.T) {
	findIndex := operators.FindIndex(
		func(val interface{}, idx int) bool {
			return val == "D"
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E").Pipe(findIndex),
			rx.Just("A", "B", "C").Pipe(findIndex),
			rx.Concat(rx.Just("A", "B", "C", "D", "E"), rx.Throw(ErrTest)).Pipe(findIndex),
			rx.Concat(rx.Just("A", "B", "C"), rx.Throw(ErrTest)).Pipe(findIndex),
		},
		[][]interface{}{
			{3, rx.Complete},
			{rx.Complete},
			{3, rx.Complete},
			{ErrTest},
		},
	)
}