package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_IsEmpty(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B").Pipe(operators.IsEmpty()),
			Just("A").Pipe(operators.IsEmpty()),
			Empty().Pipe(operators.IsEmpty()),
			Throw(errTest).Pipe(operators.IsEmpty()),
		},
		false, Complete,
		false, Complete,
		true, Complete,
		errTest,
	)
}
