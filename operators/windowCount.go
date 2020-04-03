package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// A WindowCountConfigure is a configure for WindowCount.
type WindowCountConfigure struct {
	WindowSize       int
	StartWindowEvery int
}

// Use creates an Operator from this configure.
func (configure WindowCountConfigure) Use() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return windowCountObservable{source, configure}.Subscribe
	}
}

type windowCountObservable struct {
	Source rx.Observable
	WindowCountConfigure
}

func (obs windowCountObservable) Subscribe(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
	var (
		windows    []rx.Subject
		windowSize int
	)

	if obs.StartWindowEvery == 0 {
		obs.StartWindowEvery = obs.WindowSize
	}

	window := rx.NewSubject()
	windows = append(windows, window)
	sink.Next(window.Observable)

	return obs.Source.Subscribe(ctx, func(t rx.Notification) {
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
				windows[n-1] = rx.Subject{}
				windows = windows[:n-1]
				window.Complete()
				windowSize = obs.WindowSize - obs.StartWindowEvery
				if windowSize < 0 {
					window := rx.NewSubject()
					windows = append(windows, window)
					sink.Next(window.Observable)
				}
			}

			if obs.StartWindowEvery <= obs.WindowSize {
				if windowSize%obs.StartWindowEvery == 0 {
					window := rx.NewSubject()
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
func WindowCount(windowSize int) rx.Operator {
	return WindowCountConfigure{WindowSize: windowSize}.Use()
}