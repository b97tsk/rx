package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestInterval(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			CombineLatest(
				Interval(step(1)).Pipe(operators.Take(3)),
				Interval(0).Pipe(operators.Take(3)),
			).Pipe(toString),
			CombineLatest(
				Interval(step(2)).Pipe(operators.Take(3)),
				Timer(step(1), step(2)).Pipe(operators.Take(3)),
			).Pipe(toString),
		},
		"[0 2]", "[1 2]", "[2 2]", xComplete,
		"[0 0]", "[0 1]", "[1 1]", "[1 2]", "[2 2]", xComplete,
	)
}
