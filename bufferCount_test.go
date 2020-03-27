package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_BufferCount(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(operators.BufferCount(2), toString),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(operators.BufferCount(3), toString),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(BufferCountConfigure{3, 1}.Use(), toString),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(BufferCountConfigure{3, 2}.Use(), toString),
			Just("A", "B", "C", "D", "E", "F", "G").Pipe(BufferCountConfigure{3, 4}.Use(), toString),
		},
		[][]interface{}{
			{"[A B]", "[C D]", "[E F]", "[G]", Complete},
			{"[A B C]", "[D E F]", "[G]", Complete},
			{"[A B C]", "[B C D]", "[C D E]", "[D E F]", "[E F G]", "[F G]", "[G]", Complete},
			{"[A B C]", "[C D E]", "[E F G]", "[G]", Complete},
			{"[A B C]", "[E F G]", Complete},
		},
	)
}
