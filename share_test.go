package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestShare(t *testing.T) {
	t.Parallel()

	ctx := rx.NewBackgroundContext()

	NewTestSuite[int](t).Case(
		func() rx.Observable[int] {
			ob := rx.Pipe3(
				rx.Ticker(Step(3)),
				rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
				rx.Take[int](4),
				rx.Share[int](ctx),
			)
			return rx.Merge(
				ob,
				rx.Pipe1(ob, DelaySubscription[int](4)),
				rx.Pipe1(ob, DelaySubscription[int](8)),
				rx.Pipe1(ob, DelaySubscription[int](13)),
			)
		}(),
		0, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, ErrComplete,
	).Case(
		func() rx.Observable[int] {
			ob := rx.Pipe3(
				rx.Ticker(Step(3)),
				rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
				rx.Share[int](ctx),
				rx.Take[int](4),
			)
			return rx.Merge(
				ob,
				rx.Pipe1(ob, DelaySubscription[int](4)),
				rx.Pipe1(ob, DelaySubscription[int](8)),
				rx.Pipe1(ob, DelaySubscription[int](19)),
			)
		}(),
		0, 1, 1, 2, 2, 2, 3, 3, 3, 4, 4, 5, 0, 1, 2, 3, ErrComplete,
	).Case(
		func() rx.Observable[int] {
			ob := rx.Pipe3(
				rx.Ticker(Step(3)),
				rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
				rx.Take[int](4),
				rx.Share[int](ctx).WithConnector(
					func() rx.Subject[int] {
						return rx.MulticastBuffer[int](1)
					},
				),
			)
			return rx.Merge(
				ob,
				rx.Pipe1(ob, DelaySubscription[int](4)),
				rx.Pipe1(ob, DelaySubscription[int](8)),
				rx.Pipe1(ob, DelaySubscription[int](13)),
			)
		}(),
		0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 0, 1, 2, 3, ErrComplete,
	).Case(
		func() rx.Observable[int] {
			ob := rx.Pipe3(
				rx.Ticker(Step(3)),
				rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
				rx.Share[int](ctx).WithConnector(
					func() rx.Subject[int] {
						return rx.MulticastBuffer[int](1)
					},
				),
				rx.Take[int](4),
			)
			return rx.Merge(
				ob,
				rx.Pipe1(ob, DelaySubscription[int](4)),
				rx.Pipe1(ob, DelaySubscription[int](8)),
				rx.Pipe1(ob, DelaySubscription[int](16)),
			)
		}(),
		0, 0, 1, 1, 1, 2, 2, 2, 3, 3, 3, 4, 0, 1, 2, 3, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Oops[int]("should not happen"),
			rx.Share[int](ctx).WithConnector(
				func() rx.Subject[int] {
					return rx.Subject[int]{
						Observable: rx.Throw[int](ErrTest),
						Observer:   rx.Noop[int],
					}
				},
			),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Oops[int]("should not happen"),
			rx.Share[int](ctx).WithConnector(
				func() rx.Subject[int] {
					return rx.Subject[int]{
						Observable: rx.Oops[int](ErrTest),
						Observer:   rx.Noop[int],
					}
				},
			),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Oops[int]("should not happen"),
			rx.Share[int](ctx).WithConnector(func() rx.Subject[int] { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
