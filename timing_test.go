package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestTicker(t *testing.T) {
	Subscribe(
		t,
		rx.Ticker(Step(1)).Pipe(
			operators.Map(
				func(val interface{}, idx int) interface{} {
					return idx
				},
			),
			operators.Take(3),
		),
		0, 1, 2, rx.Complete,
	)
}

func TestTimer(t *testing.T) {
	Subscribe(
		t,
		rx.Timer(Step(1)).Pipe(operators.MapTo(42)),
		42, rx.Complete,
	)
}
