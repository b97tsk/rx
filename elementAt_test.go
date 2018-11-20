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
			Concat(Range(1, 9), Throw(xErrTest)).Pipe(findFifth),
			Concat(Range(1, 5), Throw(xErrTest)).Pipe(findFifth),
		},
		5, xComplete,
		ErrOutOfRange,
		5, xComplete,
		xErrTest,
	)
}

func TestOperators_ElementAtOrDefault(t *testing.T) {
	findFifth := operators.ElementAtOrDefault(4, 404)
	subscribe(
		t,
		[]Observable{
			Range(1, 9).Pipe(findFifth),
			Range(1, 5).Pipe(findFifth),
			Concat(Range(1, 9), Throw(xErrTest)).Pipe(findFifth),
			Concat(Range(1, 5), Throw(xErrTest)).Pipe(findFifth),
		},
		5, xComplete,
		404, xComplete,
		5, xComplete,
		xErrTest,
	)
}
