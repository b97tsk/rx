package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Catch(t *testing.T) {
	op := operators.Catch(
		func(error) Observable {
			return Just("D", "E")
		},
	)
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "C").Pipe(op),
			Concat(Just("A", "B", "C"), Throw(errTest)).Pipe(op),
		},
		"A", "B", "C", Complete,
		"A", "B", "C", "D", "E", Complete,
	)
}
