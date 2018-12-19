package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_GroupBy(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Just("A", "B", "B", "A", "C", "C", "D", "A").Pipe(
				addLatencyToValue(0, 1),
				operators.GroupBy(
					func(val interface{}) interface{} {
						return val
					},
					func() Subject {
						return NewReplaySubject(0, 0).Subject
					},
				),
				operators.MergeMap(
					func(val interface{}, idx int) Observable {
						group := val.(GroupedObservable)
						delay := step(int([]rune(group.Key.(string))[0] - 'A'))
						return group.Pipe(
							operators.Count(),
							operators.Map(
								func(val interface{}, idx int) interface{} {
									return []interface{}{group.Key, val}
								},
							),
							operators.Delay(delay), // for ordered output
						)
					},
				),
				toString,
			),
		},
		"[A 3]", "[B 2]", "[C 2]", "[D 1]", xComplete,
	)
}
