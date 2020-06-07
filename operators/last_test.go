package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
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
			{1, rx.Completed},
			{2, rx.Completed},
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
			{404, rx.Completed},
			{ErrTest},
			{1, rx.Completed},
			{2, rx.Completed},
			{ErrTest},
			{ErrTest},
		},
	)
}
