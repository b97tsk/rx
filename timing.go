package rx

import (
	"time"

	"github.com/b97tsk/rx/internal/timerpool"
)

// Ticker creates an Observable that emits [time.Time] values
// every specified interval of time.
func Ticker(d time.Duration) Observable[time.Time] {
	if d <= 0 {
		return Oops[time.Time]("Ticker: d <= 0")
	}

	return func(c Context, sink Observer[time.Time]) {
		tk := time.NewTicker(d)

		c.Go(func() {
			defer tk.Stop()

			done := c.Done()

			for {
				select {
				case <-done:
					sink.Error(c.Err())
					return
				case t := <-tk.C:
					Try1(sink, Next(t), func() { sink.Error(ErrOops) })
				}
			}
		})
	}
}

// Timer creates an Observable that emits a [time.Time] value
// after a particular time span has passed, and then completes.
func Timer(d time.Duration) Observable[time.Time] {
	return func(c Context, sink Observer[time.Time]) {
		tm := timerpool.Get(d)

		c.Go(func() {
			select {
			case <-c.Done():
				timerpool.Put(tm)
				sink.Error(c.Err())
			case t := <-tm.C:
				timerpool.PutExpired(tm)
				Try1(sink, Next(t), func() { sink.Error(ErrOops) })
				sink.Complete()
			}
		})
	}
}
