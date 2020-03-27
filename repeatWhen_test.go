package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_RepeatWhen(t *testing.T) {
	var (
		repeatOnce  = operators.Take(0)
		repeatTwice = operators.Take(1)
	)

	subscribeN(
		t,
		[]Observable{
			Range(1, 4).Pipe(operators.RepeatWhen(repeatOnce)),
			Range(1, 4).Pipe(operators.RepeatWhen(repeatTwice)),
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.RepeatWhen(repeatOnce)),
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.RepeatWhen(repeatTwice)),
		},
		[][]interface{}{
			{1, 2, 3, Complete},
			{1, 2, 3, 1, 2, 3, Complete},
			{1, 2, 3, errTest},
			{1, 2, 3, errTest},
		},
	)
}
