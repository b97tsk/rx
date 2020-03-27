package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_SubscribeOn(t *testing.T) {
	subscribe(
		t,
		Merge(
			Just("A", "B").Pipe(operators.SubscribeOn(step(1))),
			Just("C", "D").Pipe(operators.SubscribeOn(step(2))),
			Just("E", "F").Pipe(operators.SubscribeOn(step(3))),
		),
		"A", "B", "C", "D", "E", "F", Complete,
	)
}
