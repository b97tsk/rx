package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestRepeat(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 4).Pipe(operators.Repeat(0)),
			rx.Range(1, 4).Pipe(operators.Repeat(1)),
			rx.Range(1, 4).Pipe(operators.Repeat(2)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.Repeat(0)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.Repeat(1)),
			rx.Concat(rx.Range(1, 4), rx.Throw(ErrTest)).Pipe(operators.Repeat(2)),
		},
		[][]interface{}{
			{rx.Completed},
			{1, 2, 3, rx.Completed},
			{1, 2, 3, 1, 2, 3, rx.Completed},
			{rx.Completed},
			{1, 2, 3, ErrTest},
			{1, 2, 3, ErrTest},
		},
	)
}
