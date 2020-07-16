package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
	"github.com/b97tsk/rx/subject"
)

func TestShare(t *testing.T) {
	t.Run("#1", shareTest1)
	t.Run("#2", shareTest2)
	t.Run("#3", shareTest3)
	t.Run("#4", shareTest4)
}

func shareTest1(t *testing.T) {
	obs := rx.Ticker(Step(3)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
		operators.Take(4),
		operators.Share(
			func() rx.Double {
				return subject.NewSubject().Double
			},
		),
	)
	Subscribe(
		t,
		rx.Merge(
			obs,
			obs.Pipe(DelaySubscription(4)),
			obs.Pipe(DelaySubscription(8)),
			obs.Pipe(DelaySubscription(13)),
		),
		0, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, Completed,
	)
}

func shareTest2(t *testing.T) {
	obs := rx.Ticker(Step(3)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
		operators.Share(
			func() rx.Double {
				return subject.NewSubject().Double
			},
		),
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
		0, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 5, 0, 1, 2, 3, Completed,
	)
}

func shareTest3(t *testing.T) {
	obs := rx.Ticker(Step(3)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
		operators.Take(4),
		operators.Share(
			func() rx.Double {
				return subject.NewReplaySubject(1).Double
			},
		),
	)
	Subscribe(
		t,
		rx.Merge(
			obs,
			obs.Pipe(DelaySubscription(4)),
			obs.Pipe(DelaySubscription(8)),
			obs.Pipe(DelaySubscription(13)),
		),
		0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, Completed,
	)
}

func shareTest4(t *testing.T) {
	obs := rx.Ticker(Step(3)).Pipe(
		operators.Map(
			func(val interface{}, idx int) interface{} {
				return idx
			},
		),
		operators.Share(
			func() rx.Double {
				return subject.NewReplaySubject(1).Double
			},
		),
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
		0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 0, 1, 2, 3, Completed,
	)
}
