package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestRetryWhen(t *testing.T) {
	retryNever := operators.Take(0)
	retryOnce := operators.Take(1)
	retryTwice := operators.Take(2)

	NewTestSuite(t).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.RetryWhen(retryNever),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.RetryWhen(retryOnce),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Just("A", "B", "C").Pipe(
			operators.RetryWhen(retryTwice),
		),
		"A", "B", "C", Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.RetryWhen(retryNever),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.RetryWhen(retryOnce),
		),
		"A", "B", "C", "A", "B", "C", ErrTest,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.RetryWhen(retryTwice),
		),
		"A", "B", "C", "A", "B", "C", "A", "B", "C", ErrTest,
	).Case(
		rx.Concat(
			rx.Just("A", "B", "C"),
			rx.Throw(ErrTest),
		).Pipe(
			operators.RetryWhen(
				func(rx.Observable) rx.Observable {
					return rx.Throw(ErrTest)
				},
			),
		),
		"A", "B", "C", ErrTest,
	).TestAll()
}
