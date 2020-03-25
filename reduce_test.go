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
			Throw(errTest).Pipe(operators.Reduce(sum)),
		},
		6, Complete,
		42, Complete,
		Complete,
		21, Complete,
		42, Complete,
		Complete,
		errTest,
	)
}

func TestOperators_Fold(t *testing.T) {
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
			Range(1, 7).Pipe(operators.Fold(-1, max)),
			Just(42).Pipe(operators.Fold(-1, max)),
			Empty().Pipe(operators.Fold(-1, max)),
			Range(1, 7).Pipe(operators.Fold(-1, sum)),
			Just(42).Pipe(operators.Fold(-1, sum)),
			Empty().Pipe(operators.Fold(-1, sum)),
			Throw(errTest).Pipe(operators.Fold(-1, sum)),
		},
		6, Complete,
		42, Complete,
		-1, Complete,
		20, Complete,
		41, Complete,
		-1, Complete,
		errTest,
	)
}
