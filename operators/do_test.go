package operators_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Do(t *testing.T) {
	n := 0
	op := operators.Do(func(Notification) { n++ })
	obs := Defer(func() Observable { return Just(n) })
	subscribeN(
		t,
		[]Observable{
			Concat(Empty().Pipe(op), obs),
			Concat(Just("A").Pipe(op), obs),
			Concat(Just("A", "B").Pipe(op), obs),
			Concat(Concat(Just("A", "B"), Throw(errTest)).Pipe(op), obs),
			obs,
		},
		[][]interface{}{
			{1, Complete},
			{"A", 3, Complete},
			{"A", "B", 6, Complete},
			{"A", "B", errTest},
			{9, Complete},
		},
	)
}

func TestOperators_DoOnNext(t *testing.T) {
	n := 0
	op := operators.DoOnNext(func(interface{}) { n++ })
	obs := Defer(func() Observable { return Just(n) })
	subscribeN(
		t,
		[]Observable{
			Concat(Empty().Pipe(op), obs),
			Concat(Just("A").Pipe(op), obs),
			Concat(Just("A", "B").Pipe(op), obs),
			Concat(Concat(Just("A", "B"), Throw(errTest)).Pipe(op), obs),
			obs,
		},
		[][]interface{}{
			{0, Complete},
			{"A", 1, Complete},
			{"A", "B", 3, Complete},
			{"A", "B", errTest},
			{5, Complete},
		},
	)
}

func TestOperators_DoOnError(t *testing.T) {
	n := 0
	op := operators.DoOnError(func(error) { n++ })
	obs := Defer(func() Observable { return Just(n) })
	subscribeN(
		t,
		[]Observable{
			Concat(Empty().Pipe(op), obs),
			Concat(Just("A").Pipe(op), obs),
			Concat(Just("A", "B").Pipe(op), obs),
			Concat(Concat(Just("A", "B"), Throw(errTest)).Pipe(op), obs),
			obs,
		},
		[][]interface{}{
			{0, Complete},
			{"A", 0, Complete},
			{"A", "B", 0, Complete},
			{"A", "B", errTest},
			{1, Complete},
		},
	)
}

func TestOperators_DoOnComplete(t *testing.T) {
	n := 0
	op := operators.DoOnComplete(func() { n++ })
	obs := Defer(func() Observable { return Just(n) })
	subscribeN(
		t,
		[]Observable{
			Concat(Empty().Pipe(op), obs),
			Concat(Just("A").Pipe(op), obs),
			Concat(Just("A", "B").Pipe(op), obs),
			Concat(Concat(Just("A", "B"), Throw(errTest)).Pipe(op), obs),
			obs,
		},
		[][]interface{}{
			{1, Complete},
			{"A", 2, Complete},
			{"A", "B", 3, Complete},
			{"A", "B", errTest},
			{3, Complete},
		},
	)
}

func TestOperators_DoAtLast(t *testing.T) {
	n := 0
	op := operators.DoAtLast(func(Notification) { n++ })
	obs := Defer(func() Observable { return Just(n) })
	subscribeN(
		t,
		[]Observable{
			Concat(Empty().Pipe(op), obs),
			Concat(Just("A").Pipe(op), obs),
			Concat(Just("A", "B").Pipe(op), obs),
			Concat(Concat(Just("A", "B"), Throw(errTest)).Pipe(op), obs),
			obs,
		},
		[][]interface{}{
			{1, Complete},
			{"A", 2, Complete},
			{"A", "B", 3, Complete},
			{"A", "B", errTest},
			{4, Complete},
		},
	)
}
