package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestCount(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(operators.Count()),
			rx.Range(1, 9).Pipe(operators.Count()),
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(operators.Count()),
		},
		[][]interface{}{
			{0, rx.Complete},
			{8, rx.Complete},
			{ErrTest},
		},
	)
}
