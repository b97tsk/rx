package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Map(t *testing.T) {
	op := operators.Map(
		func(val interface{}, idx int) interface{} {
			return val.(int) * 2
		},
	)
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(op),
			Range(1, 5).Pipe(op),
			Concat(Range(1, 5), Throw(errTest)).Pipe(op),
		},
		Complete,
		2, 4, 6, 8, Complete,
		2, 4, 6, 8, errTest,
	)
}

func TestOperators_MapTo(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.MapTo(42)),
			Just("A", "B", "C").Pipe(operators.MapTo(42)),
			Concat(Just("A", "B", "C"), Throw(errTest)).Pipe(operators.MapTo(42)),
		},
		Complete,
		42, 42, 42, Complete,
		42, 42, 42, errTest,
	)
}
