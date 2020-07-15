package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestGroupBy(t *testing.T) {
	Subscribe(
		t,
		rx.Just("A", "B", "B", "A", "C", "C", "D", "A").Pipe(
			AddLatencyToValues(0, 1),
			operators.GroupBy(
				func(val interface{}) interface{} {
					return val
				},
				func() rx.Double {
					return rx.NewReplaySubject(0).Double
				},
			),
			operators.MergeMap(
				func(val interface{}, idx int) rx.Observable {
					group := val.(rx.GroupedObservable)
					delay := Step(int([]rune(group.Key.(string))[0] - 'A'))
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
			ToString(),
		),
		"[A 3]", "[B 2]", "[C 2]", "[D 1]", Completed,
	)
}
