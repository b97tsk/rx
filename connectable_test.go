package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Share(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
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
	})
	t.Run("#2", func(t *testing.T) {
		obs := Interval(step(3)).Pipe(
			operators.Share(),
			operators.Take(4),
		)
		subscribe(
			t,
			[]Observable{
				Merge(
					obs,
					obs.Pipe(delaySubscription(4)),
					obs.Pipe(delaySubscription(8)),
					obs.Pipe(delaySubscription(20)),
				),
			},
			0, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 5, 0, 1, 2, 3, xComplete,
		)
	})
}
