package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestTakeWhile(t *testing.T) {
	takeLessThan5 := operators.TakeWhile(
		func(val interface{}, idx int) bool {
			return val.(int) < 5
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(takeLessThan5),
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(takeLessThan5),
			rx.Concat(rx.Range(1, 5), rx.Throw(ErrTest)).Pipe(takeLessThan5),
		},
		[][]interface{}{
			{1, 2, 3, 4, rx.Completed},
			{1, 2, 3, 4, rx.Completed},
			{1, 2, 3, 4, ErrTest},
		},
	)
}
