package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSome(t *testing.T) {
	someGreaterThan4 := operators.Some(
		func(val interface{}, idx int) bool {
			return val.(int) > 4
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 9).Pipe(someGreaterThan4),
			rx.Range(1, 5).Pipe(someGreaterThan4),
			rx.Empty().Pipe(someGreaterThan4),
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(someGreaterThan4),
			rx.Concat(rx.Range(1, 5), rx.Throw(ErrTest)).Pipe(someGreaterThan4),
		},
		[][]interface{}{
			{true, rx.Completed},
			{false, rx.Completed},
			{false, rx.Completed},
			{true, rx.Completed},
			{ErrTest},
		},
	)
}
