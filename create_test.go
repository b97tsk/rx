package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestCreate(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Create(
				func(ctx context.Context, sink rx.Observer) {
					sink.Next("A")
					sink.Next("B")
					sink.Next("C")
					sink.Complete()
				},
			),
			rx.Create(
				func(ctx context.Context, sink rx.Observer) {
					sink.Next("A")
					sink.Next("B")
					sink.Next("C")
					sink.Error(ErrTest)
				},
			),
		},
		[][]interface{}{
			{"A", "B", "C", rx.Completed},
			{"A", "B", "C", ErrTest},
		},
	)
}
