package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestSkip(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 7).Pipe(operators.Skip(0)),
			rx.Range(1, 7).Pipe(operators.Skip(3)),
			rx.Range(1, 3).Pipe(operators.Skip(3)),
			rx.Range(1, 1).Pipe(operators.Skip(3)),
		},
		[][]interface{}{
			{1, 2, 3, 4, 5, 6, Completed},
			{4, 5, 6, Completed},
			{Completed},
			{Completed},
		},
	)

	SubscribeN(
		t,
		[]rx.Observable{
			rx.Concat(rx.Range(1, 7), rx.Throw(ErrTest)).Pipe(operators.Skip(0)),
			rx.Concat(rx.Range(1, 7), rx.Throw(ErrTest)).Pipe(operators.Skip(3)),
			rx.Concat(rx.Range(1, 3), rx.Throw(ErrTest)).Pipe(operators.Skip(3)),
			rx.Concat(rx.Range(1, 1), rx.Throw(ErrTest)).Pipe(operators.Skip(3)),
		},
		[][]interface{}{
			{1, 2, 3, 4, 5, 6, ErrTest},
			{4, 5, 6, ErrTest},
			{ErrTest},
			{ErrTest},
		},
	)
}
