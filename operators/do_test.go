package operators_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDo(t *testing.T) {
	n := 0

	do := operators.Do(func(rx.Notification) { n++ })

	obs := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink.Next(n)
			sink.Complete()
		},
	)

	NewTestSuite(t).Case(
		rx.Concat(
			rx.Empty().Pipe(do),
			obs,
		),
		1, Completed,
	).Case(
		rx.Concat(
			rx.Just("A").Pipe(do),
			obs,
		),
		"A", 3, Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B").Pipe(do),
			obs,
		),
		"A", "B", 6, Completed,
	).Case(
		rx.Concat(
			rx.Concat(
				rx.Just("A", "B"),
				rx.Throw(ErrTest),
			).Pipe(do),
			obs,
		),
		"A", "B", ErrTest,
	).Case(
		obs,
		9, Completed,
	).TestAll()
}

func TestDoOnNext(t *testing.T) {
	n := 0

	do := operators.DoOnNext(func(interface{}) { n++ })

	obs := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink.Next(n)
			sink.Complete()
		},
	)

	NewTestSuite(t).Case(
		rx.Concat(
			rx.Empty().Pipe(do),
			obs,
		),
		0, Completed,
	).Case(
		rx.Concat(
			rx.Just("A").Pipe(do),
			obs,
		),
		"A", 1, Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B").Pipe(do),
			obs,
		),
		"A", "B", 3, Completed,
	).Case(
		rx.Concat(
			rx.Concat(
				rx.Just("A", "B"),
				rx.Throw(ErrTest),
			).Pipe(do),
			obs,
		),
		"A", "B", ErrTest,
	).Case(
		obs,
		5, Completed,
	).TestAll()
}

func TestDoOnError(t *testing.T) {
	n := 0

	do := operators.DoOnError(func(error) { n++ })

	obs := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink.Next(n)
			sink.Complete()
		},
	)

	NewTestSuite(t).Case(
		rx.Concat(
			rx.Empty().Pipe(do),
			obs,
		),
		0, Completed,
	).Case(
		rx.Concat(
			rx.Just("A").Pipe(do),
			obs,
		),
		"A", 0, Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B").Pipe(do),
			obs,
		),
		"A", "B", 0, Completed,
	).Case(
		rx.Concat(
			rx.Concat(
				rx.Just("A", "B"),
				rx.Throw(ErrTest),
			).Pipe(do),
			obs,
		),
		"A", "B", ErrTest,
	).Case(
		obs,
		1, Completed,
	).TestAll()
}

func TestDoOnComplete(t *testing.T) {
	n := 0

	do := operators.DoOnComplete(func() { n++ })

	obs := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink.Next(n)
			sink.Complete()
		},
	)

	NewTestSuite(t).Case(
		rx.Concat(
			rx.Empty().Pipe(do),
			obs,
		),
		1, Completed,
	).Case(
		rx.Concat(
			rx.Just("A").Pipe(do),
			obs,
		),
		"A", 2, Completed,
	).Case(
		rx.Concat(
			rx.Just("A", "B").Pipe(do),
			obs,
		),
		"A", "B", 3, Completed,
	).Case(
		rx.Concat(
			rx.Concat(
				rx.Just("A", "B"),
				rx.Throw(ErrTest),
			).Pipe(do),
			obs,
		),
		"A", "B", ErrTest,
	).Case(
		obs,
		3, Completed,
	).TestAll()
}
