package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_ElementAt(t *testing.T) {
	findFifth := operators.ElementAt(4)
	subscribe(
		t,
		[]Observable{
			Range(1, 9).Pipe(findFifth),
			Range(1, 5).Pipe(findFifth),
			Concat(Range(1, 9), Throw(errTest)).Pipe(findFifth),
			Concat(Range(1, 5), Throw(errTest)).Pipe(findFifth),
		},
		5, Complete,
		ErrOutOfRange,
		5, Complete,
		errTest,
	)
}

func TestOperators_ElementAtOrDefault(t *testing.T) {
	findFifth := operators.ElementAtOrDefault(4, 404)
	subscribe(
		t,
		[]Observable{
			Range(1, 9).Pipe(findFifth),
			Range(1, 5).Pipe(findFifth),
			Concat(Range(1, 9), Throw(errTest)).Pipe(findFifth),
			Concat(Range(1, 5), Throw(errTest)).Pipe(findFifth),
		},
		5, Complete,
		404, Complete,
		5, Complete,
		errTest,
	)
}
