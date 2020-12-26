package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestAudit(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Timer(Step(3))
			}),
		),
		"B", "D", "E", Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			DelaySubscription(1),
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Empty()
			}),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			DelaySubscription(1),
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Empty().Pipe(DelaySubscription(2))
			}),
		),
		Completed,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			DelaySubscription(1),
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Throw(ErrTest)
			}),
		),
		ErrTest,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Throw(ErrTest).Pipe(DelaySubscription(1))
			}),
		),
		ErrTest,
	).Case(
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.AuditTime(Step(3)),
		),
		"B", "D", "E", Completed,
	).Case(
		rx.Empty().Pipe(
			operators.AuditTime(Step(3)),
		),
		Completed,
	).Case(
		rx.Throw(ErrTest).Pipe(
			operators.AuditTime(Step(3)),
		),
		ErrTest,
	).TestAll()
}
