package operators_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestProjectToObservable(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just(
			rx.Empty(),
		).Pipe(
			operators.ConcatAll(),
		),
		Completed,
	).Case(
		rx.Just(
			func(ctx context.Context, sink rx.Observer) {
				sink.Error(ErrTest)
			},
		).Pipe(
			operators.ConcatAll(),
		),
		ErrTest,
	).Case(
		rx.Just(
			nil,
		).Pipe(
			operators.ConcatAll(),
		),
		rx.ErrNotObservable,
	).TestAll()
}
