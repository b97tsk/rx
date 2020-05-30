package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestWindowToggle(t *testing.T) {
	toSlice := func(val interface{}, idx int) rx.Observable {
		if obs, ok := val.(rx.Observable); ok {
			return obs.Pipe(
				operators.ToSlice(),
			)
		}
		return rx.Throw(rx.ErrNotObservable)
	}
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowToggle(
					rx.Ticker(Step(2)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowToggle(
					rx.Ticker(Step(2)),
					func(interface{}) rx.Observable { return rx.Timer(Step(4)) },
				),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowToggle(
					rx.Ticker(Step(4)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Concat(rx.Just("A", "B", "C", "D", "E", "F", "G"), rx.Throw(ErrTest)).Pipe(
				AddLatencyToNotifications(1, 2),
				operators.WindowToggle(
					rx.Ticker(Step(4)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
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
						if val.(int) > 1 {
							return rx.Throw(ErrTest)
						}
						return rx.Timer(Step(2))
					},
				),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowToggle(
					rx.Ticker(Step(4)).Pipe(operators.Take(2)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				operators.MergeMap(toSlice),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.WindowToggle(
					rx.Concat(rx.Ticker(Step(4)).Pipe(operators.Take(2)), rx.Throw(ErrTest)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				operators.MergeMap(toSlice),
				ToString(),
			),
		},
		[][]interface{}{
			{"[B]", "[C]", "[D]", "[E]", "[F]", "[G]", rx.Completed},
			{"[B C]", "[C D]", "[D E]", "[E F]", "[F G]", "[G]", rx.Completed},
			{"[C]", "[E]", "[G]", rx.Completed},
			{"[C]", "[E]", "[G]", ErrTest},
			{"[C]", "[E]", ErrTest},
			{"[C]", "[E]", rx.Completed},
			{"[C]", ErrTest},
		},
	)
}
