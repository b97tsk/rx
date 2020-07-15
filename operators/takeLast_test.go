package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestTakeLast(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 9).Pipe(operators.TakeLast(0)),
			rx.Range(1, 9).Pipe(operators.TakeLast(3)),
			rx.Range(1, 3).Pipe(operators.TakeLast(3)),
			rx.Range(1, 1).Pipe(operators.TakeLast(3)),
		},
		[][]interface{}{
			{Completed},
			{6, 7, 8, Completed},
			{1, 2, Completed},
			{Completed},
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(operators.TakeLast(0)),
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(operators.TakeLast(3)),
			rx.Concat(rx.Range(1, 3), rx.Throw(ErrTest)).Pipe(operators.TakeLast(3)),
			rx.Concat(rx.Range(1, 1), rx.Throw(ErrTest)).Pipe(operators.TakeLast(3)),
		},
		[][]interface{}{
			{Completed},
			{ErrTest},
			{ErrTest},
			{ErrTest},
		},
	)
}
