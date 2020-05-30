package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/testing"
)

func TestCombineLatest(t *testing.T) {
	Subscribe(
		t,
		rx.CombineLatest(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		).Pipe(ToString()),
		"[A C E]", "[A C F]", "[A D F]", "[B D F]", rx.Completed,
	)
}
