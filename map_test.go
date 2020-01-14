package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Map(t *testing.T) {
	op := operators.Map(
		func(val interface{}, idx int) interface{} {
			return val.(int) * 2
		},
	)
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(op),
			Range(1, 5).Pipe(op),
			Concat(Range(1, 5), Throw(xErrTest)).Pipe(op),
		},
		xComplete,
		2, 4, 6, 8, xComplete,
		2, 4, 6, 8, xErrTest,
	)
}

func TestOperators_MapTo(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.MapTo(42)),
			Just("A", "B", "C").Pipe(operators.MapTo(42)),
			Concat(Just("A", "B", "C"), Throw(xErrTest)).Pipe(operators.MapTo(42)),
		},
		xComplete,
		42, 42, 42, xComplete,
		42, 42, 42, xErrTest,
	)
}
