package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCongest(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe4(
			rx.Range(1, 10),
			AddLatencyToValues[int](1, 1),
			rx.Congest[int](0),
			rx.Congest[int](3),
			AddLatencyToValues[int](3, 4),
		),
		1, 2, 3, 4, 5, 6, 7, 8, 9, ErrComplete,
	)
}
