package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestWindowCount(t *testing.T) {
	toSlice := func(val interface{}, idx int) (rx.Observable, error) {
		if obs, ok := val.(rx.Observable); ok {
			return obs.Pipe(operators.ToSlice()), nil
		}
		return nil, rx.ErrNotObservable
	}
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowCount(2),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowCount(3),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowCountConfigure{3, 1}.Use(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowCountConfigure{3, 2}.Use(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowCountConfigure{3, 4}.Use(),
				operators.MergeMap(toSlice),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A B]", "[C D]", "[E F]", "[G]", rx.Completed},
			{"[A B C]", "[D E F]", "[G]", rx.Completed},
			{"[A B C]", "[B C D]", "[C D E]", "[D E F]", "[E F G]", "[F G]", "[G]", "[]", rx.Completed},
			{"[A B C]", "[C D E]", "[E F G]", "[G]", rx.Completed},
			{"[A B C]", "[E F G]", "[]", rx.Completed},
		},
	)
}
