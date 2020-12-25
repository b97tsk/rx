package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestWindowToggle(t *testing.T) {
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
				operators.WindowToggle(
					rx.Ticker(Step(2)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowToggle(
					rx.Ticker(Step(2)),
					func(interface{}) rx.Observable { return rx.Timer(Step(4)) },
				),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowToggle(
					rx.Ticker(Step(4)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Concat(rx.Just("A", "B", "C", "D", "E"), rx.Throw(ErrTest)).Pipe(
				AddLatencyToNotifications(1, 2),
				operators.WindowToggle(
					rx.Ticker(Step(4)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowToggle(
					rx.Ticker(Step(4)).Pipe(
						operators.Map(
							func(val interface{}, idx int) interface{} {
								return idx
							},
						),
					),
					func(val interface{}) rx.Observable {
						if val.(int) > 0 {
							return rx.Empty()
						}
						return rx.Timer(Step(2))
					},
				),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowToggle(
					rx.Ticker(Step(4)).Pipe(
						operators.Map(
							func(val interface{}, idx int) interface{} {
								return idx
							},
						),
					),
					func(val interface{}) rx.Observable {
						if val.(int) > 0 {
							return rx.Throw(ErrTest)
						}
						return rx.Timer(Step(2))
					},
				),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowToggle(
					rx.Ticker(Step(4)).Pipe(operators.Take(1)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowToggle(
					rx.Concat(rx.Ticker(Step(4)).Pipe(operators.Take(1)), rx.Throw(ErrTest)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				operators.MergeMap(toSlice, -1),
				ToString(),
			),
		},
		[][]interface{}{
			{"[B]", "[C]", "[D]", "[E]", Completed},
			{"[B C]", "[C D]", "[D E]", "[E]", Completed},
			{"[C]", "[E]", Completed},
			{"[C]", "[E]", ErrTest},
			{"[C]", "[E]", Completed},
			{"[C]", ErrTest},
			{"[C]", Completed},
			{ErrTest},
		},
	)
}
