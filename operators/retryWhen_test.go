package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestRetryWhen(t *testing.T) {
	var (
		retryNever = operators.Take(0)
		retryOnce  = operators.Take(1)
		retryTwice = operators.Take(2)
	)

	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 4).Pipe(operators.RetryWhen(retryNever)),
			rx.Range(1, 4).Pipe(operators.RetryWhen(retryOnce)),
			rx.Range(1, 4).Pipe(operators.RetryWhen(retryTwice)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.RetryWhen(retryNever)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.RetryWhen(retryOnce)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.RetryWhen(retryTwice)),
		},
		[][]interface{}{
			{1, 2, 3, rx.Completed},
			{1, 2, 3, rx.Completed},
			{1, 2, 3, rx.Completed},
			{1, 2, 3, ErrTest},
			{1, 2, 3, 1, 2, 3, ErrTest},
			{1, 2, 3, 1, 2, 3, 1, 2, 3, ErrTest},
		},
	)
}
