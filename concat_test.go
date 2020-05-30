package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/testing"
)

func TestConcat(t *testing.T) {
	Subscribe(
		t,
		rx.Concat(
			rx.Just("A", "B").Pipe(AddLatencyToValues(3, 5)),
			rx.Just("C", "D").Pipe(AddLatencyToValues(2, 4)),
			rx.Just("E", "F").Pipe(AddLatencyToValues(1, 3)),
		),
		"A", "B", "C", "D", "E", "F", rx.Completed,
	)
}
