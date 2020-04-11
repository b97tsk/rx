package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestShare(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
		obs := rx.Interval(Step(3)).Pipe(
			operators.Take(4),
			operators.Share(rx.NewSubject),
		)
		Subscribe(
			t,
			rx.Merge(
				obs,
				obs.Pipe(DelaySubscription(4)),
				obs.Pipe(DelaySubscription(8)),
				obs.Pipe(DelaySubscription(13)),
			),
			0, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, rx.Complete,
		)
	})
	t.Run("#2", func(t *testing.T) {
		obs := rx.Interval(Step(3)).Pipe(
			operators.Share(rx.NewSubject),
			operators.Take(4),
		)
		Subscribe(
			t,
			rx.Merge(
				obs,
				obs.Pipe(DelaySubscription(4)),
				obs.Pipe(DelaySubscription(8)),
				obs.Pipe(DelaySubscription(19)),
			),
			0, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 5, 0, 1, 2, 3, rx.Complete,
		)
	})
}

func TestShareReplay(t *testing.T) {
	t.Run("#1", func(t *testing.T) {
		obs := rx.Interval(Step(3)).Pipe(
			operators.Take(4),
			operators.ShareReplay(1, 0),
		)
		Subscribe(
			t,
			rx.Merge(
				obs,
				obs.Pipe(DelaySubscription(4)),
				obs.Pipe(DelaySubscription(8)),
				obs.Pipe(DelaySubscription(13)),
			),
			0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, rx.Complete,
		)
	})
	t.Run("#2", func(t *testing.T) {
		obs := rx.Interval(Step(3)).Pipe(
			operators.ShareReplay(1, 0),
			operators.Take(4),
		)
		Subscribe(
			t,
			rx.Merge(
				obs,
				obs.Pipe(DelaySubscription(4)),
				obs.Pipe(DelaySubscription(8)),
				obs.Pipe(DelaySubscription(16)),
			),
			0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 0, 1, 2, 3, rx.Complete,
		)
	})
}
