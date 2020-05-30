package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
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
			{rx.Completed},
			{6, 7, 8, rx.Completed},
			{1, 2, rx.Completed},
			{rx.Completed},
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
			{rx.Completed},
			{ErrTest},
			{ErrTest},
			{ErrTest},
		},
	)
}
