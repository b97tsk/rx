package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestDo(t *testing.T) {
	t.Parallel()

	doTest(t, func(f func()) rx.Operator[int, int] {
		return rx.Do(func(rx.Notification[int]) { f() })
	}, 1, 3, 6, 9)
}

func TestDoOnNext(t *testing.T) {
	t.Parallel()

	doTest(t, func(f func()) rx.Operator[int, int] {
		return rx.DoOnNext(func(int) { f() })
	}, 0, 1, 3, 5)
}

func TestDoOnComplete(t *testing.T) {
	t.Parallel()

	doTest(t, rx.DoOnComplete[int], 1, 2, 3, 3)
}

func TestDoOnError(t *testing.T) {
	t.Parallel()

	doTest(t, func(f func()) rx.Operator[int, int] {
		return rx.DoOnError[int](func(error) { f() })
	}, 0, 0, 0, 1)
}

func TestDoOnErrorOrComplete(t *testing.T) {
	t.Parallel()

	doTest(t, rx.DoOnErrorOrComplete[int], 1, 2, 3, 4)
}

func doTest(
	t *testing.T,
	op func(func()) rx.Operator[int, int],
	r1, r2, r3, r4 int,
) {
	n := 0

	do := op(func() { n++ })

	obs := rx.NewObservable(
		func(ctx context.Context, sink rx.Observer[int]) {
			sink.Next(n)
			sink.Complete()
		},
	)

	NewTestSuite[int](t).Case(
		rx.Concat(
			rx.Pipe1(rx.Empty[int](), do),
			obs,
		),
		r1, ErrCompleted,
	).Case(
		rx.Concat(
			rx.Pipe1(rx.Just(-1), do),
			obs,
		),
		-1, r2, ErrCompleted,
	).Case(
		rx.Concat(
			rx.Pipe1(rx.Just(-1, -2), do),
			obs,
		),
		-1, -2, r3, ErrCompleted,
	).Case(
		rx.Concat(
			rx.Pipe1(
				rx.Concat(
					rx.Just(-1, -2),
					rx.Throw[int](ErrTest),
				),
				do,
			),
			obs,
		),
		-1, -2, ErrTest,
	).Case(
		obs,
		r4, ErrCompleted,
	)
}
