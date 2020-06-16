package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestRetry(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 4).Pipe(operators.Retry(0)),
			rx.Range(1, 4).Pipe(operators.Retry(1)),
			rx.Range(1, 4).Pipe(operators.Retry(2)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.Retry(0)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.Retry(1)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.Retry(2)),
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
