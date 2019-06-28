package rx

import (
	"context"
)

// A WindowCountConfigure is a configure for WindowCount.
type WindowCountConfigure struct {
	WindowSize       int
	StartWindowEvery int
}

// MakeFunc creates an OperatorFunc from this type.
func (conf WindowCountConfigure) MakeFunc() OperatorFunc {
	return MakeFunc(windowCountOperator(conf).Call)
}

type windowCountOperator WindowCountConfigure

func (op windowCountOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		windows    []Subject
		windowSize int
	)

	if op.StartWindowEvery == 0 {
		op.StartWindowEvery = op.WindowSize
	}

	window := NewSubject()
	windows = append(windows, window)
	sink.Next(window.Observable)

	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			if windowSize < 0 {
				windowSize++
				break
			}

			for _, subject := range windows {
				t.Observe(subject.Observer)
			}

			windowSize++

			if windowSize == op.WindowSize {
				window := windows[0]
				copy(windows, windows[1:])
				n := len(windows)
				windows[n-1] = Subject{}
				windows = windows[:n-1]
				window.Complete()
				windowSize = op.WindowSize - op.StartWindowEvery
				if windowSize < 0 {
					window := NewSubject()
					windows = append(windows, window)
					sink.Next(window.Observable)
				}
			}

			if op.StartWindowEvery <= op.WindowSize {
				if windowSize%op.StartWindowEvery == 0 {
					window := NewSubject()
					windows = append(windows, window)
					sink.Next(window.Observable)
				}
			}

		default:
			for _, subject := range windows {
				t.Observe(subject.Observer)
			}
			sink(t)
		}
	})
}

// WindowCount branches out the source Observable values as a nested Observable
// with each nested Observable emitting at most windowSize values.
//
// It's like BufferCount, but emits a nested Observable instead of a slice.
func (Operators) WindowCount(windowSize int) OperatorFunc {
	return func(source Observable) Observable {
		op := windowCountOperator{WindowSize: windowSize}
		return source.Lift(op.Call)
	}
}
