package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestMergeScan(t *testing.T) {
	f := func(acc, val interface{}, idx int) rx.Observable {
		return val.(rx.Observable)
	}

	sum := func(seed, val interface{}, idx int) interface{} {
		return seed.(int) + val.(int)
	}

	NewTestSuite(t).Case(
		rx.Just(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		).Pipe(
			operators.MergeScan(f, nil, -1),
		),
		"E", "C", "A", "F", "D", "B", Completed,
	).Case(
		rx.Range(0, 9).Pipe(
			operators.MergeScan(
				func(acc, val interface{}, idx int) rx.Observable {
					return rx.Just(val)
				},
				nil,
				3,
			),
			operators.Reduce(sum),
		),
		36, Completed,
	).Case(
		rx.Empty().Pipe(
			operators.MergeScan(f, nil, -1),
		),
		Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.MergeScan(f, nil, -1),
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
		func() { operators.MergeScan(nil, nil, -1) },
		"MergeScan with nil accumulator didn't panic.",
	)
	panictest(
		func() {
			operators.MergeScan(
				func(interface{}, interface{}, int) rx.Observable {
					return rx.Throw(ErrTest)
				},
				nil,
				0,
			)
		},
		"MergeScan with zero concurrency didn't panic.",
	)
}
