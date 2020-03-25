package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_IgnoreElements(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.IgnoreElements()),
			Just("A", "B", "C").Pipe(operators.IgnoreElements()),
			Concat(Just("A", "B", "C"), Throw(errTest)).Pipe(operators.IgnoreElements()),
		},
		Complete,
		Complete,
		errTest,
	)
}
