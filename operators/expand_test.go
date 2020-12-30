package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestExpand(t *testing.T) {
	sum := func(seed, val interface{}, idx int) interface{} {
		return seed.(int) + val.(int)
	}

	NewTestSuite(t).Case(
		rx.Just(5).Pipe(
			operators.Expand(
				func(val interface{}) rx.Observable {
					i := val.(int)
					if i < 1 {
						return rx.Empty()
					}

					return rx.Just(i-1, i-1).Pipe(
						AddLatencyToValues(1, 1),
					)
				},
				3,
			),
			operators.Reduce(sum),
		),
		57, Completed,
	).Case(
		rx.Just(5).Pipe(
			operators.Expand(
				func(val interface{}) rx.Observable {
					i := val.(int)
					if i < 1 {
						return rx.Throw(ErrTest)
					}

					return rx.Just(i - 1)
				},
				-1,
			),
		),
		5, 4, 3, 2, 1, 0, ErrTest,
	).Case(
		rx.Empty().Pipe(
			operators.Expand(
				func(interface{}) rx.Observable {
					return rx.Throw(ErrTest)
				},
				-1,
			),
		),
		Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Expand(
				func(interface{}) rx.Observable {
					return rx.Throw(ErrTest)
				},
				-1,
			),
		),
		ErrTest,
	).TestAll()

	panictest := func(f func(), msg string) {
		defer func() {
			if recover() == nil {
				t.Log(msg)
				t.FailNow()
			}
		}()
		f()
	}
	panictest(
		func() { operators.Expand(nil, -1) },
		"Expand with nil project didn't panic.",
	)
	panictest(
		func() {
			operators.Expand(
				func(interface{}) rx.Observable {
					return rx.Throw(ErrTest)
				},
				0,
			)
		},
		"Expand with zero concurrency didn't panic.",
	)
}
