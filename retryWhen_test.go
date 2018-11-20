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

	subscribe(
		t,
		[]Observable{
			Range(1, 4).Pipe(operators.RetryWhen(retryNever)),
			Range(1, 4).Pipe(operators.RetryWhen(retryOnce)),
			Range(1, 4).Pipe(operators.RetryWhen(retryTwice)),
			Concat(Range(1, 4), Throw(xErrTest)).Pipe(operators.RetryWhen(retryNever)),
			Concat(Range(1, 4), Throw(xErrTest)).Pipe(operators.RetryWhen(retryOnce)),
			Concat(Range(1, 4), Throw(xErrTest)).Pipe(operators.RetryWhen(retryTwice)),
		},
		1, 2, 3, xComplete,
		1, 2, 3, xComplete,
		1, 2, 3, xComplete,
		1, 2, 3, xErrTest,
		1, 2, 3, 1, 2, 3, xErrTest,
		1, 2, 3, 1, 2, 3, 1, 2, 3, xErrTest,
	)
}
