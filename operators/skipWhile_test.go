package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestSkipWhile(t *testing.T) {
	skipLessThan5 := operators.SkipWhile(
		func(val interface{}, idx int) bool {
			return val.(int) < 5
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(skipLessThan5),
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(skipLessThan5),
			rx.Concat(rx.Range(1, 5), rx.Throw(ErrTest)).Pipe(skipLessThan5),
		},
		[][]interface{}{
			{5, 4, 3, 2, 1, rx.Completed},
			{5, 6, 7, 8, ErrTest},
			{ErrTest},
		},
	)
}
