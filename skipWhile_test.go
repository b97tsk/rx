package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_SkipWhile(t *testing.T) {
	skipLessThan5 := operators.SkipWhile(
		func(val interface{}, idx int) bool {
			return val.(int) < 5
		},
	)
	subscribe(
		t,
		[]Observable{
			Just(1, 2, 3, 4, 5, 4, 3, 2, 1).Pipe(skipLessThan5),
			Concat(Range(1, 9), Throw(errTest)).Pipe(skipLessThan5),
			Concat(Range(1, 5), Throw(errTest)).Pipe(skipLessThan5),
		},
		5, 4, 3, 2, 1, Complete,
		5, 6, 7, 8, errTest,
		errTest,
	)
}
