package operators_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestCongest(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Range(1, 9).Pipe(
			AddLatencyToValues(1, 1),
			operators.Congest(0),
			operators.Congest(3),
			AddLatencyToValues(3, 4),
		),
		1, 2, 3, 4, 5, 6, 7, 8, Completed,
	).TestAll()

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	rx.Timer(Step(2)).Pipe(
		operators.Congest(3),
	).Subscribe(ctx, func(rx.Notification) { t.Fatal("should not happen") })
}
