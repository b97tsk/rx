package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestLast(t *testing.T) {
	last := operators.Last()
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(last),
			rx.Throw(ErrTest).Pipe(last),
			rx.Just(1).Pipe(last),
			rx.Just(1, 2).Pipe(last),
			rx.Concat(rx.Just(1), rx.Throw(ErrTest)).Pipe(last),
			rx.Concat(rx.Just(1, 2), rx.Throw(ErrTest)).Pipe(last),
		},
		[][]interface{}{
			{rx.ErrEmpty},
			{ErrTest},
			{1, Completed},
			{2, Completed},
			{ErrTest},
			{ErrTest},
		},
	)
}

func TestLastOrDefault(t *testing.T) {
	last := operators.LastOrDefault(404)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(last),
			rx.Throw(ErrTest).Pipe(last),
			rx.Just(1).Pipe(last),
			rx.Just(1, 2).Pipe(last),
			rx.Concat(rx.Just(1), rx.Throw(ErrTest)).Pipe(last),
			rx.Concat(rx.Just(1, 2), rx.Throw(ErrTest)).Pipe(last),
		},
		[][]interface{}{
			{404, Completed},
			{ErrTest},
			{1, Completed},
			{2, Completed},
			{ErrTest},
			{ErrTest},
		},
	)
}
