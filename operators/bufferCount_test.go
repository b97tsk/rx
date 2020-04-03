package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestBufferCount(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(operators.BufferCount(2), ToString()),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(operators.BufferCount(3), ToString()),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(operators.BufferCountConfigure{3, 1}.Use(), ToString()),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(operators.BufferCountConfigure{3, 2}.Use(), ToString()),
			rx.Just("A", "B", "C", "D", "E", "F", "G").Pipe(operators.BufferCountConfigure{3, 4}.Use(), ToString()),
		},
		[][]interface{}{
			{"[A B]", "[C D]", "[E F]", "[G]", rx.Complete},
			{"[A B C]", "[D E F]", "[G]", rx.Complete},
			{"[A B C]", "[B C D]", "[C D E]", "[D E F]", "[E F G]", "[F G]", "[G]", rx.Complete},
			{"[A B C]", "[C D E]", "[E F G]", "[G]", rx.Complete},
			{"[A B C]", "[E F G]", rx.Complete},
		},
	)
}