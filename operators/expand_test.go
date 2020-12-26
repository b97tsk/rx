package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestExpand(t *testing.T) {
	NewTestSuite(t).Case(
		rx.Just(8).Pipe(
			operators.Expand(
				func(val interface{}) rx.Observable {
					i := val.(int)
					if i < 1 {
						return rx.Empty()
					}
					return rx.Just(i - 1)
				},
				-1,
			),
		),
		8, 7, 6, 5, 4, 3, 2, 1, 0, Completed,
	).TestAll()
}
