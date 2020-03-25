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
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.DefaultIfEmpty(42)),
		},
		42, Complete,
		1, 2, 3, Complete,
		1, 2, 3, errTest,
	)
}
