package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestWindow(t *testing.T) {
	toSlice := func(val interface{}, idx int) rx.Observable {
		if obs, ok := val.(rx.Observable); ok {
			return obs.Pipe(operators.ToSlice())
		}

		return rx.Throw(rx.ErrNotObservable)
	}

	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.Window(rx.Ticker(Step(2))),
			operators.MergeMap(toSlice, -1),
			ToString(),
		),
		"[A]", "[B]", "[C]", "[D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.Window(rx.Ticker(Step(4))),
			operators.MergeMap(toSlice, -1),
			ToString(),
		),
		"[A B]", "[C D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.Window(rx.Ticker(Step(6))),
			operators.MergeMap(toSlice, -1),
			ToString(),
		),
		"[A B C]", "[D E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.Window(rx.Ticker(Step(8))),
			operators.MergeMap(toSlice, -1),
			ToString(),
		),
		"[A B C D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.Window(rx.Throw(ErrTest)),
			operators.MergeMap(toSlice, -1),
			ToString(),
		),
		ErrTest,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.Window(rx.Ticker(Step(1))),
			operators.MergeMap(toSlice, -1),
			ToString(),
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
		func() { operators.Window(nil) },
		"Window with nil window boundaries didn't panic.",
	)
}
