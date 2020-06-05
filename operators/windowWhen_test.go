package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestWindowWhen(t *testing.T) {
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
				operators.WindowWhen(func() rx.Observable { return rx.Timer(Step(2)) }),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowWhen(func() rx.Observable { return rx.Timer(Step(4)) }),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowWhen(func() rx.Observable { return rx.Timer(Step(6)) }),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowWhen(func() rx.Observable { return rx.Timer(Step(8)) }),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowWhen(func() rx.Observable { return rx.Throw(ErrTest) }),
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
