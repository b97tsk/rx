package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestTicker(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe2(
			rx.Ticker(Step(1)),
			rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
			rx.Take[int](3),
		),
		0, 1, 2, ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Ticker(0),
			rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
			rx.Take[int](3),
		),
		rx.ErrOops, "Ticker: d <= 0",
	).Case(
		rx.Pipe2(
			rx.Ticker(Step(1)),
			rx.Scan(-1, func(i int, _ time.Time) int { return i + 1 }),
			rx.OnNext(func(int) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}

func TestTimer(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.MapTo[time.Time](42),
		),
		42, ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Timer(Step(1)),
			rx.MapTo[time.Time](42),
			rx.OnNext(func(int) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	NewTestSuite[int](t).WithContext(ctx).Case(
		rx.Pipe1(
			rx.Timer(Step(2)),
			rx.MapTo[time.Time](42),
		),
		context.DeadlineExceeded,
	)
}
