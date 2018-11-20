package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Reduce(t *testing.T) {
	max := func(seed, val interface{}, idx int) interface{} {
		if seed.(int) > val.(int) {
			return seed
		}
		return val
	}

	sum := func(seed, val interface{}, idx int) interface{} {
		return seed.(int) + val.(int)
	}

	subscribe(
		t,
		[]Observable{
			Range(1, 7).Pipe(operators.Reduce(max)),
			Just(42).Pipe(operators.Reduce(max)),
			Empty().Pipe(operators.Reduce(max)),
			Range(1, 7).Pipe(operators.Reduce(sum)),
			Just(42).Pipe(operators.Reduce(sum)),
			Empty().Pipe(operators.Reduce(sum)),
			Throw(xErrTest).Pipe(operators.Reduce(sum)),
		},
		6, xComplete,
		42, xComplete,
		xComplete,
		21, xComplete,
		42, xComplete,
		xComplete,
		xErrTest,
	)
}
