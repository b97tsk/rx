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
func (configure WindowCountConfigure) MakeFunc() OperatorFunc {
	return func(source Observable) Observable {
		return windowCountObservable{source, configure}.Subscribe
	}
}

type windowCountObservable struct {
	Source Observable
	WindowCountConfigure
}

func (obs windowCountObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var (
		windows    []Subject
		windowSize int
	)

	if obs.StartWindowEvery == 0 {
		obs.StartWindowEvery = obs.WindowSize
	}

	window := NewSubject()
	windows = append(windows, window)
	sink.Next(window.Observable)

	return obs.Source.Subscribe(ctx, func(t Notification) {
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

			if windowSize == obs.WindowSize {
				window := windows[0]
				copy(windows, windows[1:])
				n := len(windows)
				windows[n-1] = Subject{}
				windows = windows[:n-1]
				window.Complete()
				windowSize = obs.WindowSize - obs.StartWindowEvery
				if windowSize < 0 {
					window := NewSubject()
					windows = append(windows, window)
					sink.Next(window.Observable)
				}
			}

			if obs.StartWindowEvery <= obs.WindowSize {
				if windowSize%obs.StartWindowEvery == 0 {
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
	return WindowCountConfigure{WindowSize: windowSize}.MakeFunc()
}
