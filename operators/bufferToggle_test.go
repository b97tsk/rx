package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestBufferToggle(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferToggle(
					rx.Ticker(Step(2)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferToggle(
					rx.Ticker(Step(2)),
					func(interface{}) rx.Observable { return rx.Timer(Step(4)) },
				),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferToggle(
					rx.Ticker(Step(4)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				ToString(),
			),
			rx.Concat(rx.Just("A", "B", "C", "D", "E", "F", "G"), rx.Throw(ErrTest)).Pipe(
				AddLatencyToNotifications(1, 2),
				operators.BufferToggle(
					rx.Ticker(Step(4)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferToggle(
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
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferToggle(
					rx.Ticker(Step(4)).Pipe(operators.Take(2)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
				ToString(),
			),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				AddLatencyToValues(1, 2),
				operators.BufferToggle(
					rx.Concat(rx.Ticker(Step(4)).Pipe(operators.Take(2)), rx.Throw(ErrTest)),
					func(interface{}) rx.Observable { return rx.Timer(Step(2)) },
				),
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
