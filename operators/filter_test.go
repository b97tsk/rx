package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestFilter(t *testing.T) {
	filterLessThan5 := operators.Filter(
		func(val interface{}, idx int) bool {
			return val.(int) < 5
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(filterLessThan5),
			rx.Range(1, 9).Pipe(filterLessThan5),
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(filterLessThan5),
		},
		[][]interface{}{
			{1, 2, 3, 4, 4, 3, 2, 1, rx.Completed},
			{1, 2, 3, 4, rx.Completed},
			{1, 2, 3, 4, ErrTest},
		},
	)
}

func TestExclude(t *testing.T) {
	excludeLessThan5 := operators.Exclude(
		func(val interface{}, idx int) bool {
			return val.(int) < 5
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(excludeLessThan5),
			rx.Range(1, 9).Pipe(excludeLessThan5),
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(excludeLessThan5),
		},
		[][]interface{}{
			{5, rx.Completed},
			{5, 6, 7, 8, rx.Completed},
			{5, 6, 7, 8, ErrTest},
		},
	)
}
