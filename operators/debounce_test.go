package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDebounce(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C").Pipe(
				AddLatencyToValues(1, 2),
				operators.Debounce(func(interface{}) rx.Observable {
					return rx.Timer(Step(3))
				}),
			),
			rx.Just("A", "B", "C").Pipe(
				AddLatencyToValues(1, 3),
				operators.Debounce(func(interface{}) rx.Observable {
					return rx.Timer(Step(2))
				}),
			),
			rx.Just("A", "B", "C").Pipe(
				AddLatencyToValues(1, 3),
				operators.Debounce(func(interface{}) rx.Observable {
					return rx.Empty()
				}),
			),
			rx.Just("A", "B", "C").Pipe(
				AddLatencyToValues(1, 3),
				operators.Debounce(func(interface{}) rx.Observable {
					return rx.Throw(ErrTest)
				}),
			),
			rx.Throw(ErrTest).Pipe(
				operators.Debounce(func(interface{}) rx.Observable {
					return rx.Timer(Step(1))
				}),
			),
		},
		[][]interface{}{
			{"C", Completed},
			{"A", "B", "C", Completed},
			{"C", Completed},
			{ErrTest},
			{ErrTest},
		},
	)
}

func TestDebounceTime(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C").Pipe(
				AddLatencyToValues(1, 2),
				operators.DebounceTime(Step(3)),
			),
			rx.Just("A", "B", "C").Pipe(
				AddLatencyToValues(1, 3),
				operators.DebounceTime(Step(2)),
			),
		},
		[][]interface{}{
			{"C", Completed},
			{"A", "B", "C", Completed},
		},
	)
}
