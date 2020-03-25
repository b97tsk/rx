package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Every(t *testing.T) {
	everyLessThan5 := operators.Every(
		func(val interface{}, idx int) bool {
			return val.(int) < 5
		},
	)
	subscribe(
		t,
		[]Observable{
			Range(1, 9).Pipe(everyLessThan5),
			Range(1, 5).Pipe(everyLessThan5),
			Empty().Pipe(everyLessThan5),
			Concat(Range(1, 9), Throw(errTest)).Pipe(everyLessThan5),
			Concat(Range(1, 5), Throw(errTest)).Pipe(everyLessThan5),
		},
		false, Complete,
		true, Complete,
		true, Complete,
		false, Complete,
		errTest,
	)
}
