package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestWindowCount(t *testing.T) {
	toSlice := func(val interface{}, idx int) rx.Observable {
		if obs, ok := val.(rx.Observable); ok {
			return obs.Pipe(operators.ToSlice())
		}

		return rx.Throw(rx.ErrNotObservable)
	}

	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.WindowCount(2),
			operators.MergeMap(toSlice, -1),
			ToString(),
		),
		"[A B]", "[C D]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.WindowCount(3),
			operators.MergeMap(toSlice, -1),
			ToString(),
		),
		"[A B C]", "[D E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.WindowCountConfigure{
				WindowSize:       3,
				StartWindowEvery: 1,
			}.Make(),
			operators.MergeMap(toSlice, -1),
			ToString(),
		),
		"[A B C]", "[B C D]", "[C D E]", "[D E]", "[E]", "[]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.WindowCountConfigure{
				WindowSize:       3,
				StartWindowEvery: 2,
			}.Make(),
			operators.MergeMap(toSlice, -1),
			ToString(),
		),
		"[A B C]", "[C D E]", "[E]", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.WindowCountConfigure{
				WindowSize:       3,
				StartWindowEvery: 4,
			}.Make(),
			operators.MergeMap(toSlice, -1),
			ToString(),
		),
		"[A B C]", "[E]", Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.WindowCount(2),
			operators.MergeMap(toSlice, -1),
			ToString(),
		),
		ErrTest,
	).TestAll()
}
