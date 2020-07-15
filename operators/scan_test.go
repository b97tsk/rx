package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestScan(t *testing.T) {
	max := func(acc, val interface{}, idx int) interface{} {
		if acc.(int) > val.(int) {
			return acc
		}
		return val
	}

	sum := func(acc, val interface{}, idx int) interface{} {
		return acc.(int) + val.(int)
	}

	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 7).Pipe(operators.Scan(max)),
			rx.Just(42).Pipe(operators.Scan(max)),
			rx.Empty().Pipe(operators.Scan(max)),
			rx.Range(1, 7).Pipe(operators.Scan(sum)),
			rx.Just(42).Pipe(operators.Scan(sum)),
			rx.Empty().Pipe(operators.Scan(sum)),
			rx.Throw(ErrTest).Pipe(operators.Scan(sum)),
		},
		[][]interface{}{
			{1, 2, 3, 4, 5, 6, Completed},
			{42, Completed},
			{Completed},
			{1, 3, 6, 10, 15, 21, Completed},
			{42, Completed},
			{Completed},
			{ErrTest},
		},
	)
}
