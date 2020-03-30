package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Buffer(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Buffer(Interval(step(2))),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Buffer(Interval(step(4))),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Buffer(Interval(step(6))),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Buffer(Interval(step(8))),
				toString,
			),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(
				addLatencyToValue(1, 2),
				operators.Buffer(Throw(errTest)),
				toString,
			),
		},
		[][]interface{}{
			{"[A]", "[B]", "[C]", "[D]", "[E]", "[F]", Complete},
			{"[A B]", "[C D]", "[E F]", Complete},
			{"[A B C]", "[D E F]", Complete},
			{"[A B C D]", Complete},
			{errTest},
		},
	)
}
