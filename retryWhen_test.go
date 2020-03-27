package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_RetryWhen(t *testing.T) {
	var (
		retryNever = operators.Take(0)
		retryOnce  = operators.Take(1)
		retryTwice = operators.Take(2)
	)

	subscribeN(
		t,
		[]Observable{
			Range(1, 4).Pipe(operators.RetryWhen(retryNever)),
			Range(1, 4).Pipe(operators.RetryWhen(retryOnce)),
			Range(1, 4).Pipe(operators.RetryWhen(retryTwice)),
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.RetryWhen(retryNever)),
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.RetryWhen(retryOnce)),
			Concat(Range(1, 4), Throw(errTest)).Pipe(operators.RetryWhen(retryTwice)),
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
