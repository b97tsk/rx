package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestReduce(t *testing.T) {
	max := func(seed, val interface{}, idx int) interface{} {
		if seed.(int) > val.(int) {
			return seed
		}
		return val
	}
	sum := func(seed, val interface{}, idx int) interface{} {
		return seed.(int) + val.(int)
	}
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 7).Pipe(operators.Reduce(max)),
			rx.Just(42).Pipe(operators.Reduce(max)),
			rx.Empty().Pipe(operators.Reduce(max)),
			rx.Range(1, 7).Pipe(operators.Reduce(sum)),
			rx.Just(42).Pipe(operators.Reduce(sum)),
			rx.Empty().Pipe(operators.Reduce(sum)),
			rx.Throw(ErrTest).Pipe(operators.Reduce(sum)),
		},
		[][]interface{}{
			{6, Completed},
			{42, Completed},
			{Completed},
			{21, Completed},
			{42, Completed},
			{Completed},
			{ErrTest},
		},
	)
}

func TestFold(t *testing.T) {
	max := func(seed, val interface{}, idx int) interface{} {
		if seed.(int) > val.(int) {
			return seed
		}
		return val
	}
	sum := func(seed, val interface{}, idx int) interface{} {
		return seed.(int) + val.(int)
	}
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 7).Pipe(operators.Fold(-1, max)),
			rx.Just(42).Pipe(operators.Fold(-1, max)),
			rx.Empty().Pipe(operators.Fold(-1, max)),
			rx.Range(1, 7).Pipe(operators.Fold(-1, sum)),
			rx.Just(42).Pipe(operators.Fold(-1, sum)),
			rx.Empty().Pipe(operators.Fold(-1, sum)),
			rx.Throw(ErrTest).Pipe(operators.Fold(-1, sum)),
		},
		[][]interface{}{
			{6, Completed},
			{42, Completed},
			{-1, Completed},
			{20, Completed},
			{41, Completed},
			{-1, Completed},
			{ErrTest},
		},
	)
}
