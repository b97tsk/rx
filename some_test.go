package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Some(t *testing.T) {
	someGreaterThan4 := operators.Some(
		func(val interface{}, idx int) bool {
			return val.(int) > 4
		},
	)
	subscribe(
		t,
		[]Observable{
			Range(1, 9).Pipe(someGreaterThan4),
			Range(1, 5).Pipe(someGreaterThan4),
			Empty().Pipe(someGreaterThan4),
			Concat(Range(1, 9), Throw(errTest)).Pipe(someGreaterThan4),
			Concat(Range(1, 5), Throw(errTest)).Pipe(someGreaterThan4),
		},
		true, Complete,
		false, Complete,
		false, Complete,
		true, Complete,
		errTest,
	)
}
