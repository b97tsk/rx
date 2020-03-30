package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestFind(t *testing.T) {
	findFive := operators.Find(
		func(val interface{}, idx int) bool {
			return val.(int) == 5
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(findFive),
			rx.Range(1, 9).Pipe(findFive),
			rx.Range(1, 5).Pipe(findFive),
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(findFive),
			rx.Concat(rx.Range(1, 5), rx.Throw(ErrTest)).Pipe(findFive),
		},
		[][]interface{}{
			{5, rx.Complete},
			{5, rx.Complete},
			{rx.Complete},
			{5, rx.Complete},
			{ErrTest},
		},
	)
}
