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
			{1, 2, 3, 4, 4, 3, 2, 1, Completed},
			{1, 2, 3, 4, Completed},
			{1, 2, 3, 4, ErrTest},
		},
	)
}

func TestFilterMap(t *testing.T) {
	filterLessThan5 := operators.FilterMap(
		func(val interface{}, idx int) (interface{}, bool) {
			return val, val.(int) < 5
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
			{1, 2, 3, 4, 4, 3, 2, 1, Completed},
			{1, 2, 3, 4, Completed},
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
			{5, Completed},
			{5, 6, 7, 8, Completed},
			{5, 6, 7, 8, ErrTest},
		},
	)
}
