package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestAudit(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](1, 2),
			rx.Audit(
				func(string) rx.Observable[time.Time] {
					return rx.Timer(Step(3))
				},
			),
		),
		"B", "D", "E", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			DelaySubscription[string](1),
			rx.Audit(
				func(string) rx.Observable[int] {
					return rx.Empty[int]()
				},
			),
		),
		ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			DelaySubscription[string](1),
			rx.Audit(
				func(string) rx.Observable[int] {
					return rx.Pipe1(
						rx.Empty[int](),
						DelaySubscription[int](2),
					)
				},
			),
		),
		ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			DelaySubscription[string](1),
			rx.Audit(
				func(string) rx.Observable[int] {
					return rx.Throw[int](ErrTest)
				},
			),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C", "D", "E"),
			rx.Audit(
				func(string) rx.Observable[int] {
					return rx.Pipe1(
						rx.Throw[int](ErrTest),
						DelaySubscription[int](1),
					)
				},
			),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C", "D", "E"),
			AddLatencyToValues[string](1, 2),
			rx.AuditTime[string](Step(3)),
		),
		"B", "D", "E", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.AuditTime[string](Step(3)),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.AuditTime[string](Step(3)),
		),
		ErrTest,
	)
}
