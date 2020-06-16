package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDefaultIfEmpty(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(operators.DefaultIfEmpty(42)),
			rx.Range(1, 4).Pipe(operators.DefaultIfEmpty(42)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.DefaultIfEmpty(42)),
		},
		[][]interface{}{
			{42, rx.Completed},
			{1, 2, 3, rx.Completed},
			{1, 2, 3, ErrTest},
		},
	)
}
