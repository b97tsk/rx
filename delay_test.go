package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Delay(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Range(1, 5).Pipe(operators.Delay(step(3))),
			Concat(Range(1, 5), Throw(errTest)).Pipe(
				addLatencyToNotification(0, 3),
				operators.Delay(step(1)),
			),
		},
		[][]interface{}{
			{1, 2, 3, 4, Complete},
			{1, 2, 3, 4, errTest},
		},
	)
}
