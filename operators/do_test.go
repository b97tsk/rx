package operators_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestDo(t *testing.T) {
	n := 0
	op := operators.Do(func(rx.Notification) { n++ })
	obs := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink.Next(n)
			sink.Complete()
		},
	)
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
			{1, Completed},
			{"A", 3, Completed},
			{"A", "B", 6, Completed},
			{"A", "B", ErrTest},
			{9, Completed},
		},
	)
}

func TestDoOnNext(t *testing.T) {
	n := 0
	op := operators.DoOnNext(func(interface{}) { n++ })
	obs := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink.Next(n)
			sink.Complete()
		},
	)
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
			{0, Completed},
			{"A", 1, Completed},
			{"A", "B", 3, Completed},
			{"A", "B", ErrTest},
			{5, Completed},
		},
	)
}

func TestDoOnError(t *testing.T) {
	n := 0
	op := operators.DoOnError(func(error) { n++ })
	obs := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink.Next(n)
			sink.Complete()
		},
	)
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
			{0, Completed},
			{"A", 0, Completed},
			{"A", "B", 0, Completed},
			{"A", "B", ErrTest},
			{1, Completed},
		},
	)
}

func TestDoOnComplete(t *testing.T) {
	n := 0
	op := operators.DoOnComplete(func() { n++ })
	obs := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink.Next(n)
			sink.Complete()
		},
	)
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
			{1, Completed},
			{"A", 2, Completed},
			{"A", "B", 3, Completed},
			{"A", "B", ErrTest},
			{3, Completed},
		},
	)
}

func TestDoAtLast(t *testing.T) {
	n := 0
	op := operators.DoAtLast(func(error) { n++ })
	obs := rx.Observable(
		func(ctx context.Context, sink rx.Observer) {
			sink.Next(n)
			sink.Complete()
		},
	)
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
			{1, Completed},
			{"A", 2, Completed},
			{"A", "B", 3, Completed},
			{"A", "B", ErrTest},
			{4, Completed},
		},
	)
}
