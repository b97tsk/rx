package rxtest

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/b97tsk/rx"
)

var (
	ErrComplete = errors.New("complete")
	ErrTest     = errors.New("test")
)

func Step(n int) time.Duration {
	return 60 * time.Millisecond * time.Duration(n)
}

func AddLatencyToValues[T any](initialDelay, period int) rx.Operator[T, T] {
	return rx.NewOperator(
		func(source rx.Observable[T]) rx.Observable[T] {
			return rx.Zip2(
				source,
				rx.Concat(rx.Timer(Step(initialDelay)), rx.Ticker(Step(period))),
				func(v T, _ time.Time) T { return v },
			)
		},
	)
}

func AddLatencyToNotifications[T any](initialDelay, period int) rx.Operator[T, T] {
	return rx.NewOperator(
		func(source rx.Observable[T]) rx.Observable[T] {
			return rx.Pipe1(
				rx.Zip2(
					rx.Pipe1(source, rx.Materialize[T]()),
					rx.Concat(rx.Timer(Step(initialDelay)), rx.Ticker(Step(period))),
					func(n rx.Notification[T], _ time.Time) rx.Notification[T] { return n },
				),
				rx.Dematerialize[rx.Notification[T]](),
			)
		},
	)
}

func DelaySubscription[T any](n int) rx.Operator[T, T] {
	return rx.NewOperator(
		func(source rx.Observable[T]) rx.Observable[T] {
			return rx.Concat(
				rx.Pipe1(
					rx.Timer(Step(n)),
					rx.IgnoreElements[time.Time, T](),
				),
				source,
			)
		},
	)
}

func tos(v any) string {
	if _, ok := v.(error); ok {
		return fmt.Sprintf("<%v>", v)
	}

	return fmt.Sprint(v)
}

func ToString[T any]() rx.Operator[T, string] {
	return rx.Map(func(v T) string { return tos(v) })
}

type TestSuite[T any] struct {
	tb  testing.TB
	ctx context.Context
}

func NewTestSuite[T any](tb testing.TB) *TestSuite[T] {
	return &TestSuite[T]{tb, context.Background()}
}

func (s *TestSuite[T]) WithContext(ctx context.Context) *TestSuite[T] {
	return &TestSuite[T]{s.tb, ctx}
}

func (s *TestSuite[T]) Case(obs rx.Observable[T], output ...any) *TestSuite[T] {
	_ = obs.BlockingSubscribe(s.ctx, func(n rx.Notification[T]) {
		if len(output) == 0 {
			s.tb.Fail()

			switch {
			case n.HasValue:
				s.tb.Logf("want <nothing>, but got %v", tos(n.Value))
			case n.HasError:
				s.tb.Logf("want <nothing>, but got %v", tos(n.Error))
			default:
				s.tb.Logf("want <nothing>, but got %v", tos(ErrComplete))
			}

			return
		}

		wanted := output[0]
		output = output[1:]

		switch {
		case n.HasValue:
			if wanted != any(n.Value) {
				s.tb.Fail()
				s.tb.Logf("want %v, but got %v", tos(wanted), tos(n.Value))
			} else {
				s.tb.Logf("want %v", tos(wanted))
			}
		case n.HasError:
			if wanted != n.Error {
				s.tb.Fail()
				s.tb.Logf("want %v, but got %v", tos(wanted), tos(n.Error))
			} else {
				s.tb.Logf("want %v", tos(wanted))
			}
		default:
			if wanted != ErrComplete {
				s.tb.Fail()
				s.tb.Logf("want %v, but got %v", tos(wanted), tos(ErrComplete))
			} else {
				s.tb.Logf("want %v", tos(ErrComplete))
			}
		}
	})

	if len(output) > 0 {
		s.tb.Fail()

		for _, wanted := range output {
			s.tb.Logf("want %v, but got <nothing>", tos(wanted))
		}
	}

	return s
}
