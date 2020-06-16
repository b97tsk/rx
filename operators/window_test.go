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
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Window(rx.Ticker(Step(2))),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Window(rx.Ticker(Step(4))),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Window(rx.Ticker(Step(6))),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Window(rx.Ticker(Step(8))),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.Window(rx.Throw(ErrTest)),
				operators.MergeMap(toSlice),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", rx.Completed},
			{"[A B]", "[C D]", "[E F]", "[G]", rx.Completed},
			{"[A B C]", "[D E F]", "[G]", rx.Completed},
			{"[A B C D]", "[E F G]", rx.Completed},
			{ErrTest},
		},
	)
}
