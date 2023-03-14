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

func TestOnNext(t *testing.T) {
	t.Parallel()

	doTest(t, func(f func()) rx.Operator[int, int] {
		return rx.OnNext(func(int) { f() })
	}, 0, 1, 3, 5)
}

func TestOnComplete(t *testing.T) {
	t.Parallel()

	doTest(t, rx.OnComplete[int], 1, 2, 3, 3)
}

func TestOnError(t *testing.T) {
	t.Parallel()

	doTest(t, func(f func()) rx.Operator[int, int] {
		return rx.OnError[int](func(error) { f() })
	}, 0, 0, 0, 1)
}

func TestOnLastNotification(t *testing.T) {
	t.Parallel()

	doTest(t, rx.OnLastNotification[int], 1, 2, 3, 4)
}

func doTest(
	t *testing.T,
	op func(func()) rx.Operator[int, int],
	r1, r2, r3, r4 int,
) {
	n := 0

	do := op(func() { n++ })

	obs := rx.NewObservable(
		func(_ context.Context, sink rx.Observer[int]) {
			sink.Next(n)
			sink.Complete()
		},
	)

	NewTestSuite[int](t).Case(
		rx.Concat(
			rx.Pipe1(rx.Empty[int](), do),
			obs,
		),
		r1, ErrComplete,
	).Case(
		rx.Concat(
			rx.Pipe1(rx.Just(-1), do),
			obs,
		),
		-1, r2, ErrComplete,
	).Case(
		rx.Concat(
			rx.Pipe1(rx.Just(-1, -2), do),
			obs,
		),
		-1, -2, r3, ErrComplete,
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
		r4, ErrComplete,
	)
}
