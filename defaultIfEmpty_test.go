package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_DefaultIfEmpty(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.DefaultIfEmpty(42)),
			Range(1, 4).Pipe(operators.DefaultIfEmpty(42)),
			Concat(Range(1, 4), Throw(xErrTest)).Pipe(operators.DefaultIfEmpty(42)),
		},
		42, xComplete,
		1, 2, 3, xComplete,
		1, 2, 3, xErrTest,
	)
}
