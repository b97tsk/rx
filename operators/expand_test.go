package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestExpand(t *testing.T) {
	Subscribe(
		t,
		rx.Just(8).Pipe(operators.Expand(
			func(val interface{}, idx int) rx.Observable {
				i := val.(int)
				if i < 1 {
					return rx.Empty()
				}
				return rx.Just(i - 1)
			},
		)),
		8, 7, 6, 5, 4, 3, 2, 1, 0, rx.Completed,
	)
}
