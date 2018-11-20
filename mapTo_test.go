package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_MapTo(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.MapTo(42)),
			Just("A", "B", "C").Pipe(operators.MapTo(42)),
			Concat(Just("A", "B", "C"), Throw(xErrTest)).Pipe(operators.MapTo(42)),
		},
		xComplete,
		42, 42, 42, xComplete,
		42, 42, 42, xErrTest,
	)
}
