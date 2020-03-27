package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_DelayWhen(t *testing.T) {
	subscribe(
		t,
		Range(1, 5).Pipe(operators.DelayWhen(
			func(val interface{}, idx int) Observable {
				return Interval(step(val.(int)))
			},
		)),
		1, 2, 3, 4, Complete,
	)
}
