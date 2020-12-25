package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestWindowWhen(t *testing.T) {
	toSlice := func(val interface{}, idx int) rx.Observable {
		if obs, ok := val.(rx.Observable); ok {
			return obs.Pipe(operators.ToSlice())
		}

		return rx.Throw(rx.ErrNotObservable)
	}

	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowWhen(func() rx.Observable { return rx.Timer(Step(2)) }),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowWhen(func() rx.Observable { return rx.Timer(Step(4)) }),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowWhen(func() rx.Observable { return rx.Timer(Step(6)) }),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowWhen(func() rx.Observable { return rx.Timer(Step(8)) }),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowWhen(rx.Empty),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowWhen(func() rx.Observable { return rx.Throw(ErrTest) }),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Throw(ErrTest).Pipe(
				operators.WindowWhen(func() rx.Observable { return rx.Timer(Step(1)) }),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", Completed},
			{"[A B]", "[C D]", "[E]", Completed},
			{"[A B C]", "[D E]", Completed},
			{"[A B C D]", "[E]", Completed},
			{"[A B C D E]", Completed},
			{ErrTest},
			{ErrTest},
		},
	)
}
