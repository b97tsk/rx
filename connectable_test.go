package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Share(t *testing.T) {
	obs := Interval(step(3)).Pipe(
		operators.Take(4),
		operators.Share(),
	)
	subscribe(
		t,
		[]Observable{
			Merge(
				obs,
				obs.Pipe(delaySubscription(4)),
				obs.Pipe(delaySubscription(8)),
				obs.Pipe(delaySubscription(15)),
			),
		},
		0, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, xComplete,
	)
}
