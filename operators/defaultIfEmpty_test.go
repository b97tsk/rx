package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
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
			{42, rx.Complete},
			{1, 2, 3, rx.Complete},
			{1, 2, 3, ErrTest},
		},
	)
}
