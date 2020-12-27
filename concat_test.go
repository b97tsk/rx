package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestConcat(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Concat(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		),
		"A", "B", "C", "D", "E", "F", Completed,
	).Case(
		rx.Concat(),
		Completed,
	).TestAll()

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	rx.Concat(
		func(_ context.Context, sink rx.Observer) {
			_ = rx.Timer(Step(2)).BlockingSubscribe(context.Background(), sink)
		},
		func(context.Context, rx.Observer) {
			t.Fatal("should not happen")
		},
	).Subscribe(ctx, rx.Noop)
}
