package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestWindowTime(t *testing.T) {
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
				operators.WindowTime(Step(2)),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTime(Step(4)),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTime(Step(6)),
				operators.MergeMap(toSlice),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", rx.Completed},
			{"[A B]", "[C D]", "[E F]", "[G]", rx.Completed},
			{"[A B C]", "[D E F]", "[G]", rx.Completed},
		},
	)
	t.Log("----------")
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(8), 0, 0}.Use(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(8), 0, 3}.Use(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(8), 0, 2}.Use(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(8), 0, 1}.Use(),
				operators.MergeMap(toSlice),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A B C D]", "[E F G]", rx.Completed},
			{"[A B C]", "[D E F]", "[G]", rx.Completed},
			{"[A B]", "[C D]", "[E F]", "[G]", rx.Completed},
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", "[]", rx.Completed},
		},
	)
	t.Log("----------")
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(2), Step(2), 0}.Use(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(2), Step(4), 0}.Use(),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowTimeConfigure{Step(4), Step(2), 0}.Use(),
				operators.MergeMap(toSlice),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", "[G]", rx.Completed},
			{"[A]", "[C]", "[E]", "[G]", rx.Completed},
			{"[A B]", "[B C]", "[C D]", "[D E]", "[E F]", "[F G]", "[G]", rx.Completed},
		},
	)
}
