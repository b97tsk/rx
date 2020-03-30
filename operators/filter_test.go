package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Filter(t *testing.T) {
	filterLessThan5 := operators.Filter(
		func(val interface{}, idx int) bool {
			return val.(int) < 5
		},
	)
	subscribeN(
		t,
		[]Observable{
			Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(filterLessThan5),
			Range(1, 9).Pipe(filterLessThan5),
			Concat(Range(1, 9), Throw(errTest)).Pipe(filterLessThan5),
		},
		[][]interface{}{
			{1, 2, 3, 4, 4, 3, 2, 1, Complete},
			{1, 2, 3, 4, Complete},
			{1, 2, 3, 4, errTest},
		},
	)
}

func TestOperators_Exclude(t *testing.T) {
	excludeLessThan5 := operators.Exclude(
		func(val interface{}, idx int) bool {
			return val.(int) < 5
		},
	)
	subscribeN(
		t,
		[]Observable{
			Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(excludeLessThan5),
			Range(1, 9).Pipe(excludeLessThan5),
			Concat(Range(1, 9), Throw(errTest)).Pipe(excludeLessThan5),
		},
		[][]interface{}{
			{5, Complete},
			{5, 6, 7, 8, Complete},
			{5, 6, 7, 8, errTest},
		},
	)
}
