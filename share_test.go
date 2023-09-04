package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestShare1(t *testing.T) {
	t.Parallel()

	obs := rx.Pipe3(
		rx.Ticker(Step(3)),
		rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
		rx.Take[int](4),
		rx.Share[int](),
	)

	NewTestSuite[int](t).Case(
		rx.Merge(
			obs,
			rx.Pipe1(obs, DelaySubscription[int](4)),
			rx.Pipe1(obs, DelaySubscription[int](8)),
			rx.Pipe1(obs, DelaySubscription[int](13)),
		),
		0, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, ErrComplete,
	)
}

func TestShare2(t *testing.T) {
	t.Parallel()

	obs := rx.Pipe3(
		rx.Ticker(Step(3)),
		rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
		rx.Share[int](),
		rx.Take[int](4),
	)

	NewTestSuite[int](t).Case(
		rx.Merge(
			obs,
			rx.Pipe1(obs, DelaySubscription[int](4)),
			rx.Pipe1(obs, DelaySubscription[int](8)),
			rx.Pipe1(obs, DelaySubscription[int](19)),
		),
		0, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 5, 0, 1, 2, 3, ErrComplete,
	)
}

func TestShare3(t *testing.T) {
	t.Parallel()

	obs := rx.Pipe3(
		rx.Ticker(Step(3)),
		rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
		rx.Take[int](4),
		rx.Share[int]().WithConnector(
			func() rx.Subject[int] {
				return rx.MulticastReplay[int](&rx.ReplayConfig{BufferSize: 1})
			},
		),
	)

	NewTestSuite[int](t).Case(
		rx.Merge(
			obs,
			rx.Pipe1(obs, DelaySubscription[int](4)),
			rx.Pipe1(obs, DelaySubscription[int](8)),
			rx.Pipe1(obs, DelaySubscription[int](13)),
		),
		0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, ErrComplete,
	)
}

func TestShare4(t *testing.T) {
	t.Parallel()

	obs := rx.Pipe3(
		rx.Ticker(Step(3)),
		rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
		rx.Share[int]().WithConnector(
			func() rx.Subject[int] {
				return rx.MulticastReplay[int](&rx.ReplayConfig{BufferSize: 1})
			},
		),
		rx.Take[int](4),
	)

	NewTestSuite[int](t).Case(
		rx.Merge(
			obs,
			rx.Pipe1(obs, DelaySubscription[int](4)),
			rx.Pipe1(obs, DelaySubscription[int](8)),
			rx.Pipe1(obs, DelaySubscription[int](16)),
		),
		0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 0, 1, 2, 3, ErrComplete,
	)
}

func TestShare5(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			func(context.Context, rx.Observer[int]) {
				panic("should not happen")
			},
			rx.Share[int]().WithConnector(
				func() rx.Subject[int] {
					return rx.Subject[int]{
						Observable: rx.Throw[int](ErrTest),
						Observer:   rx.Noop[int],
					}
				},
			),
		),
		ErrTest,
	)
}
