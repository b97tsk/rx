package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestAudit(t *testing.T) {
	Subscribe(
		t,
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Timer(Step(3))
			}),
		),
		"B", "D", "E", Completed,
	)
	Subscribe(
		t,
		rx.Just("A", "B", "C", "D", "E").Pipe(
			DelaySubscription(1),
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Empty()
			}),
		),
		Completed,
	)
	Subscribe(
		t,
		rx.Just("A", "B", "C", "D", "E").Pipe(
			DelaySubscription(1),
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Empty().Pipe(DelaySubscription(2))
			}),
		),
		Completed,
	)
	Subscribe(
		t,
		rx.Just("A", "B", "C", "D", "E").Pipe(
			DelaySubscription(1),
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Throw(ErrTest)
			}),
		),
		ErrTest,
	)
	Subscribe(
		t,
		rx.Just("A", "B", "C", "D", "E").Pipe(
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Throw(ErrTest).Pipe(DelaySubscription(1))
			}),
		),
		ErrTest,
	)
}

func TestAuditTime(t *testing.T) {
	audit := operators.AuditTime(Step(3))
	Subscribe(
		t,
		rx.Just("A", "B", "C", "D", "E").Pipe(
			AddLatencyToValues(1, 2),
			audit,
		),
		"B", "D", "E", Completed,
	)
	Subscribe(t, rx.Empty().Pipe(audit), Completed)
	Subscribe(t, rx.Throw(ErrTest).Pipe(audit), ErrTest)
}
