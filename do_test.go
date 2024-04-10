package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestDo(t *testing.T) {
	t.Parallel()

	doTest(t, func(f func()) rx.Operator[int, int] {
		return rx.Do(func(rx.Notification[int]) { f() })
	}, 1, 3, 6, 9)
	doTest(t, func(f func()) rx.Operator[int, int] {
		return rx.DoOnNext(func(int) { f() })
	}, 0, 1, 3, 5)
	doTest(t, func(f func()) rx.Operator[int, int] {
		return rx.DoOnError[int](func(error) { f() })
	}, 0, 0, 0, 1)
	doTest(t, rx.DoOnComplete[int], 1, 2, 3, 3)
	doTest(t, rx.DoOnTermination[int], 1, 2, 3, 4)

	NewTestSuite[string](t).Case(
		rx.Pipe1(rx.Just("A"), rx.Do(func(n rx.Notification[string]) { panic(ErrTest) })),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(rx.Just("A"), rx.DoOnNext(func(string) { panic(ErrTest) })),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(rx.Throw[string](ErrTest), rx.DoOnError[string](func(err error) { panic(err) })),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(rx.Empty[string](), rx.DoOnComplete[string](func() { panic(ErrTest) })),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(rx.Throw[string](ErrTest), rx.DoOnTermination[string](func() { panic(ErrTest) })),
		rx.ErrOops, ErrTest,
	)
}

func doTest(
	t *testing.T,
	op func(func()) rx.Operator[int, int],
	r1, r2, r3, r4 int,
) {
	n := 0

	do := op(func() { n++ })

	obs := rx.NewObservable(
		func(_ rx.Context, o rx.Observer[int]) {
			o.Next(n)
			o.Complete()
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
