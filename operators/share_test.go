package operators_test

import (
	"context"
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
			Merge(
				obs,
				obs.Pipe(delaySubscription(4)),
				obs.Pipe(delaySubscription(8)),
				obs.Pipe(delaySubscription(13)),
			),
			0, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, Complete,
		)
	})
	t.Run("#2", func(t *testing.T) {
		obs := Interval(step(3)).Pipe(
			operators.Share(),
			operators.Take(4),
		)
		subscribe(
			t,
			Merge(
				obs,
				obs.Pipe(delaySubscription(4)),
				obs.Pipe(delaySubscription(8)),
				obs.Pipe(delaySubscription(19)),
			),
			0, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 5, 0, 1, 2, 3, Complete,
		)
	})
}

func TestOperators_ShareReplay(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
		obs := Interval(step(3)).Pipe(
			operators.Take(4),
			operators.ShareReplay(1, 0),
		)
		subscribe(
			t,
			Merge(
				obs,
				obs.Pipe(delaySubscription(4)),
				obs.Pipe(delaySubscription(8)),
				obs.Pipe(delaySubscription(13)),
			),
			0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, Complete,
		)
	})
	t.Run("#2", func(t *testing.T) {
		obs := Interval(step(3)).Pipe(
			operators.ShareReplay(1, 0),
			operators.Take(4),
		)
		subscribe(
			t,
			Merge(
				obs,
				obs.Pipe(delaySubscription(4)),
				obs.Pipe(delaySubscription(8)),
				obs.Pipe(delaySubscription(16)),
			),
			0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 0, 1, 2, 3, Complete,
		)
	})
}
