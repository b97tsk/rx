package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestDo(t *testing.T) {
	n := 0
	op := operators.Do(func(rx.Notification) { n++ })
	obs := rx.Defer(func() rx.Observable { return rx.Just(n) })
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Concat(rx.Empty().Pipe(op), obs),
			rx.Concat(rx.Just("A").Pipe(op), obs),
			rx.Concat(rx.Just("A", "B").Pipe(op), obs),
			rx.Concat(rx.Concat(rx.Just("A", "B"), rx.Throw(ErrTest)).Pipe(op), obs),
			obs,
		},
		[][]interface{}{
			{1, rx.Complete},
			{"A", 3, rx.Complete},
			{"A", "B", 6, rx.Complete},
			{"A", "B", ErrTest},
			{9, rx.Complete},
		},
	)
}

func TestDoOnNext(t *testing.T) {
	n := 0
	op := operators.DoOnNext(func(interface{}) { n++ })
	obs := rx.Defer(func() rx.Observable { return rx.Just(n) })
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Concat(rx.Empty().Pipe(op), obs),
			rx.Concat(rx.Just("A").Pipe(op), obs),
			rx.Concat(rx.Just("A", "B").Pipe(op), obs),
			rx.Concat(rx.Concat(rx.Just("A", "B"), rx.Throw(ErrTest)).Pipe(op), obs),
			obs,
		},
		[][]interface{}{
			{0, rx.Complete},
			{"A", 1, rx.Complete},
			{"A", "B", 3, rx.Complete},
			{"A", "B", ErrTest},
			{5, rx.Complete},
		},
	)
}

func TestDoOnError(t *testing.T) {
	n := 0
	op := operators.DoOnError(func(error) { n++ })
	obs := rx.Defer(func() rx.Observable { return rx.Just(n) })
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Concat(rx.Empty().Pipe(op), obs),
			rx.Concat(rx.Just("A").Pipe(op), obs),
			rx.Concat(rx.Just("A", "B").Pipe(op), obs),
			rx.Concat(rx.Concat(rx.Just("A", "B"), rx.Throw(ErrTest)).Pipe(op), obs),
			obs,
		},
		[][]interface{}{
			{0, rx.Complete},
			{"A", 0, rx.Complete},
			{"A", "B", 0, rx.Complete},
			{"A", "B", ErrTest},
			{1, rx.Complete},
		},
	)
}

func TestDoOnComplete(t *testing.T) {
	n := 0
	op := operators.DoOnComplete(func() { n++ })
	obs := rx.Defer(func() rx.Observable { return rx.Just(n) })
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Concat(rx.Empty().Pipe(op), obs),
			rx.Concat(rx.Just("A").Pipe(op), obs),
			rx.Concat(rx.Just("A", "B").Pipe(op), obs),
			rx.Concat(rx.Concat(rx.Just("A", "B"), rx.Throw(ErrTest)).Pipe(op), obs),
			obs,
		},
		[][]interface{}{
			{1, rx.Complete},
			{"A", 2, rx.Complete},
			{"A", "B", 3, rx.Complete},
			{"A", "B", ErrTest},
			{3, rx.Complete},
		},
	)
}

func TestDoAtLast(t *testing.T) {
	n := 0
	op := operators.DoAtLast(func(rx.Notification) { n++ })
	obs := rx.Defer(func() rx.Observable { return rx.Just(n) })
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Concat(rx.Empty().Pipe(op), obs),
			rx.Concat(rx.Just("A").Pipe(op), obs),
			rx.Concat(rx.Just("A", "B").Pipe(op), obs),
			rx.Concat(rx.Concat(rx.Just("A", "B"), rx.Throw(ErrTest)).Pipe(op), obs),
			obs,
		},
		[][]interface{}{
			{1, rx.Complete},
			{"A", 2, rx.Complete},
			{"A", "B", 3, rx.Complete},
			{"A", "B", ErrTest},
			{4, rx.Complete},
		},
	)
}
