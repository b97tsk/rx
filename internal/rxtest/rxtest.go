package rxtest

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
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
		return fmt.Sprintf("(%v)", v)
	}

	return fmt.Sprint(v)
}

func ToString[T any]() rx.Operator[T, string] {
	return rx.Map(func(v T) string { return tos(v) })
}

type TestSuite[T any] struct {
	tb testing.TB
	c  rx.Context
}

func NewTestSuite[T any](tb testing.TB) *TestSuite[T] {
	return &TestSuite[T]{tb, rx.NewBackgroundContext()}
}

func (s *TestSuite[T]) WithContext(c rx.Context) *TestSuite[T] {
	return &TestSuite[T]{s.tb, c}
}

func (s *TestSuite[T]) Case(ob rx.Observable[T], output ...any) *TestSuite[T] {
	var wg sync.WaitGroup

	var p atomic.Pointer[any]

	f := func(v any) { p.Store(&v) }

	_ = ob.BlockingSubscribe(s.c.WithWaitGroup(&wg).WithPanicHandler(f), func(n rx.Notification[T]) {
		if len(output) == 0 {
			switch n.Kind {
			case rx.KindNext:
				s.tb.Errorf("want (nothing), but got %v", tos(n.Value))
			case rx.KindComplete:
				s.tb.Errorf("want (nothing), but got %v", tos(ErrComplete))
			case rx.KindError, rx.KindStop:
				s.tb.Errorf("want (nothing), but got %v", tos(n.Error))
			}

			return
		}

		wanted := output[0]
		output = output[1:]

		switch n.Kind {
		case rx.KindNext:
			if wanted == any(n.Value) {
				s.tb.Logf("want %v", tos(wanted))
			} else {
				s.tb.Errorf("want %v, but got %v", tos(wanted), tos(n.Value))
			}
		case rx.KindComplete:
			if wanted == ErrComplete {
				s.tb.Logf("want %v", tos(ErrComplete))
			} else {
				s.tb.Errorf("want %v, but got %v", tos(wanted), tos(ErrComplete))
			}
		case rx.KindError, rx.KindStop:
			if wanted == n.Error {
				s.tb.Logf("want %v", tos(wanted))
			} else {
				s.tb.Errorf("want %v, but got %v", tos(wanted), tos(n.Error))
			}
		}
	})

	c := make(chan struct{})

	go func() {
		wg.Wait()
		close(c)
	}()

	tm := time.NewTimer(5 * time.Second)

	select {
	case <-c:
		tm.Stop()
	case <-tm.C:
		s.tb.Error("timeout waiting for WaitGroup")
		return s
	}

	if p := p.Load(); p != nil {
		v := *p

		if len(output) == 0 {
			s.tb.Errorf("want (nothing), but got %v", tos(v))
			return s
		}

		wanted := output[0]
		output = output[1:]

		if wanted == v {
			s.tb.Logf("want %v", tos(wanted))
		} else {
			s.tb.Errorf("want %v, but got %v", tos(wanted), tos(v))
		}
	}

	if len(output) != 0 {
		for _, wanted := range output {
			s.tb.Logf("want %v, but got (nothing)", tos(wanted))
		}
		s.tb.Fail()
	}

	return s
}
