package operators_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func testDo(t *testing.T, op func(func()) rx.Operator, r1, r2, r3, r4 int) {
	n := 0

	do := op(func() { n++ })

	obs := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink.Next(n)
			sink.Complete()
		},
	)

	NewTestSuite(t).Case(
		rx.Concat(
			rx.Empty().Pipe(do),
			obs,
		),
		r1, Completed,
	).Case(
		rx.Concat(
			rx.Just("A").Pipe(do),
			obs,
		),
		"A", r2, Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B").Pipe(do),
			obs,
		),
		"A", "B", r3, Completed,
	).Case(
		rx.Concat(
			rx.Concat(
				rx.Just("A", "B"),
				rx.Throw(ErrTest),
			).Pipe(do),
			obs,
		),
		"A", "B", ErrTest,
	).Case(
		obs,
		r4, Completed,
	).TestAll()
}

func TestDo(t *testing.T) {
	t.Run("Do", func(t *testing.T) {
		testDo(
			t,
			func(f func()) rx.Operator {
				return operators.Do(func(rx.Notification) { f() })
			},
			1, 3, 6, 9,
		)
	})
	t.Run("DoOnNext", func(t *testing.T) {
		testDo(
			t,
			func(f func()) rx.Operator {
				return operators.DoOnNext(func(interface{}) { f() })
			},
			0, 1, 3, 5,
		)
	})
	t.Run("DoOnError", func(t *testing.T) {
		testDo(
			t,
			func(f func()) rx.Operator {
				return operators.DoOnError(func(error) { f() })
			},
			0, 0, 0, 1,
		)
	})
	t.Run("DoOnComplete", func(t *testing.T) {
		testDo(t, operators.DoOnComplete, 1, 2, 3, 3)
	})
	t.Run("DoOnErrorOrComplete", func(t *testing.T) {
		testDo(t, operators.DoOnErrorOrComplete, 1, 2, 3, 4)
	})
}
