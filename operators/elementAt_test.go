package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestElementAt(t *testing.T) {
	findFifth := operators.ElementAt(4)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 9).Pipe(findFifth),
			rx.Range(1, 5).Pipe(findFifth),
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(findFifth),
			rx.Concat(rx.Range(1, 5), rx.Throw(ErrTest)).Pipe(findFifth),
		},
		[][]interface{}{
			{5, rx.Complete},
			{rx.ErrOutOfRange},
			{5, rx.Complete},
			{ErrTest},
		},
	)
}

func TestElementAtOrDefault(t *testing.T) {
	findFifth := operators.ElementAtOrDefault(4, 404)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Range(1, 9).Pipe(findFifth),
			rx.Range(1, 5).Pipe(findFifth),
			rx.Concat(rx.Range(1, 9), rx.Throw(ErrTest)).Pipe(findFifth),
			rx.Concat(rx.Range(1, 5), rx.Throw(ErrTest)).Pipe(findFifth),
		},
		[][]interface{}{
			{5, rx.Complete},
			{404, rx.Complete},
			{5, rx.Complete},
			{ErrTest},
		},
	)
}
