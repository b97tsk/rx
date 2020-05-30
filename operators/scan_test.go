package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
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
			{1, 2, 3, 4, 5, 6, rx.Completed},
			{42, rx.Completed},
			{rx.Completed},
			{1, 3, 6, 10, 15, 21, rx.Completed},
			{42, rx.Completed},
			{rx.Completed},
			{ErrTest},
		},
	)
}
