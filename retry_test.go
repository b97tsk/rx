package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Retry(t *testing.T) {
	subscribeN(
		t,
		[]Observable{
			Range(1, 4).Pipe(operators.Retry(0)),
			Range(1, 4).Pipe(operators.Retry(1)),
			Range(1, 4).Pipe(operators.Retry(2)),
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.Retry(0)),
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.Retry(1)),
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.Retry(2)),
		},
		[][]interface{}{
			{1, 2, 3, Complete},
			{1, 2, 3, Complete},
			{1, 2, 3, Complete},
			{1, 2, 3, errTest},
			{1, 2, 3, 1, 2, 3, errTest},
			{1, 2, 3, 1, 2, 3, 1, 2, 3, errTest},
		},
	)
}
