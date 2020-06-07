package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestFirst(t *testing.T) {
	first := operators.First()
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(first),
			rx.Throw(ErrTest).Pipe(first),
			rx.Just(1).Pipe(first),
			rx.Just(1, 2).Pipe(first),
			rx.Concat(rx.Just(1), rx.Throw(ErrTest)).Pipe(first),
			rx.Concat(rx.Just(1, 2), rx.Throw(ErrTest)).Pipe(first),
		},
		[][]interface{}{
			{rx.ErrEmpty},
			{ErrTest},
			{1, rx.Completed},
			{1, rx.Completed},
			{1, rx.Completed},
			{1, rx.Completed},
		},
	)
}

func TestFirstOrDefault(t *testing.T) {
	first := operators.FirstOrDefault(404)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(first),
			rx.Throw(ErrTest).Pipe(first),
			rx.Just(1).Pipe(first),
			rx.Just(1, 2).Pipe(first),
			rx.Concat(rx.Just(1), rx.Throw(ErrTest)).Pipe(first),
			rx.Concat(rx.Just(1, 2), rx.Throw(ErrTest)).Pipe(first),
		},
		[][]interface{}{
			{404, rx.Completed},
			{ErrTest},
			{1, rx.Completed},
			{1, rx.Completed},
			{1, rx.Completed},
			{1, rx.Completed},
		},
	)
}
