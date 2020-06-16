package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDelayWhen(t *testing.T) {
	Subscribe(
		t,
		rx.Range(1, 5).Pipe(operators.DelayWhen(
			func(val interface{}, idx int) rx.Observable {
				return rx.Timer(Step(val.(int)))
			},
		)),
		1, 2, 3, 4, rx.Completed,
	)
}
