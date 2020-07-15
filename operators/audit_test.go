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
		rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
			AddLatencyToValues(1, 2),
			operators.Audit(func(interface{}) rx.Observable {
				return rx.Timer(Step(3))
			}),
		),
		"B", "D", "F", Completed,
	)
}

func TestAuditTime(t *testing.T) {
	Subscribe(
		t,
		rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(
			AddLatencyToValues(1, 2),
			operators.AuditTime(Step(3)),
		),
		"B", "D", "F", Completed,
	)
}
