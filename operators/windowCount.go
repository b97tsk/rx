package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// WindowCount branches out the source values as a nested Observable with
// each nested Observable emitting at most windowSize values.
//
// It's like BufferCount, but emits a nested Observable instead of a slice.
func WindowCount(windowSize int) rx.Operator {
	return WindowCountConfigure{WindowSize: windowSize}.Make()
}

// A WindowCountConfigure is a configure for WindowCount.
type WindowCountConfigure struct {
	WindowSize       int
	StartWindowEvery int
	WindowFactory    rx.DoubleFactory
}

// Make creates an Operator from this configure.
func (configure WindowCountConfigure) Make() rx.Operator {
	if configure.WindowSize <= 0 {
		panic("WindowCount: WindowSize is negative or zero")
	}

	if configure.StartWindowEvery < 0 {
		panic("WindowCount: StartWindowEvery is negative")
	}

	if configure.StartWindowEvery == 0 {
		configure.StartWindowEvery = configure.WindowSize
	}

	if configure.WindowFactory == nil {
		configure.WindowFactory = rx.Multicast
	}

	return func(source rx.Observable) rx.Observable {
		return windowCountObservable{source, configure}.Subscribe
	}
}

type windowCountObservable struct {
	Source rx.Observable
	WindowCountConfigure
}

func (obs windowCountObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var (
		windows    []rx.Observer
		windowSize int
	)

	window := obs.WindowFactory()
	windows = append(windows, window.Observer)
	sink.Next(window.Observable)

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			if windowSize < 0 {
				windowSize++
				break
			}

			for _, window := range windows {
				window.Sink(t)
			}

			windowSize++

			if windowSize == obs.WindowSize {
				window := windows[0]

				copy(windows, windows[1:])

				n := len(windows)
				windows[n-1] = nil
				windows = windows[:n-1]

				window.Complete()

				windowSize = obs.WindowSize - obs.StartWindowEvery

				if windowSize < 0 {
					window := obs.WindowFactory()
					windows = append(windows, window.Observer)
					sink.Next(window.Observable)
				}
			}

			if obs.StartWindowEvery <= obs.WindowSize {
				if windowSize%obs.StartWindowEvery == 0 {
					window := obs.WindowFactory()
					windows = append(windows, window.Observer)
					sink.Next(window.Observable)
				}
			}

		case t.HasError:
			fallthrough

		default:
			for _, window := range windows {
				window.Sink(t)
			}

			sink(t)
		}
	})
}
