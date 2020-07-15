package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
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
			{Completed},
			{1, 2, 3, Completed},
			{1, 2, 3, 1, 2, 3, Completed},
			{Completed},
			{1, 2, 3, ErrTest},
			{1, 2, 3, ErrTest},
		},
	)
}
