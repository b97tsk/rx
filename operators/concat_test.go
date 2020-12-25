package operators_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestConcat(t *testing.T) {
	Subscribe(
		t,
		rx.Just(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		).Pipe(operators.ConcatAll()),
		"A", "B", "C", "D", "E", "F", Completed,
	)
	Subscribe(
		t,
		rx.Timer(Step(1)).Pipe(
			operators.ConcatMapTo(rx.Just("A")),
		),
		"A", Completed,
	)
	Subscribe(
		t,
		rx.Throw(ErrTest).Pipe(
			operators.ConcatAll(),
		),
		ErrTest,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	_ = rx.Just(
		func(_ context.Context, sink rx.Observer) {
			_ = rx.Timer(Step(2)).BlockingSubscribe(context.Background(), sink)
		},
		func(context.Context, rx.Observer) {
			t.Fatal("should not happen")
		},
	).Pipe(
		operators.ConcatAll(),
	).BlockingSubscribe(ctx, rx.Noop)
}
