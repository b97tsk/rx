package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestCongest(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Range(1, 9).Pipe(
			AddLatencyToValues(1, 1),
			operators.Congest(3),
			AddLatencyToValues(3, 4),
		),
		1, 2, 3, 4, 5, 6, 7, 8, Completed,
	).TestAll()
}
