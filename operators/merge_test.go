package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestMerge(t *testing.T) {
	sum := func(seed, val interface{}, idx int) interface{} {
		return seed.(int) + val.(int)
	}

	NewTestSuite(t).Case(
		rx.Just(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		).Pipe(
			operators.MergeAll(-1),
		),
		"E", "C", "A", "F", "D", "B", Completed,
	).Case(
		rx.Range(0, 9).Pipe(
			operators.MergeMap(
				func(val interface{}, idx int) rx.Observable {
					return rx.Just(val)
				},
				3,
			),
			operators.Reduce(sum),
		),
		36, Completed,
	).Case(
		rx.Timer(Step(1)).Pipe(
			operators.MergeMapTo(rx.Just("A"), -1),
		),
		"A", Completed,
	).Case(
		rx.Empty().Pipe(
			operators.MergeAll(-1),
		),
		Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.MergeAll(-1),
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
		func() { operators.MergeMap(nil, -1) },
		"MergeMap with nil project didn't panic.",
	)
	panictest(
		func() {
			operators.MergeMap(
				func(interface{}, int) rx.Observable {
					return rx.Throw(ErrTest)
				},
				0,
			)
		},
		"MergeMap with zero concurrency didn't panic.",
	)
}
