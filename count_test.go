package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Count(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.Count()),
			Range(1, 9).Pipe(operators.Count()),
			Concat(Range(1, 9), Throw(xErrTest)).Pipe(operators.Count()),
		},
		0, xComplete,
		8, xComplete,
		xErrTest,
	)
}
