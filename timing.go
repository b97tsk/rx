package rx

import "time"

// Ticker creates an Observable that emits [time.Time] values
// every specified interval of time.
func Ticker(d time.Duration) Observable[time.Time] {
	if d <= 0 {
		return Oops[time.Time]("Ticker: d <= 0")
	}

	return func(c Context, o Observer[time.Time]) {
		tk := time.NewTicker(d)
		c.Go(func() {
			defer tk.Stop()
			done := c.Done()
			for {
				select {
				case <-done:
					o.Error(c.Cause())
					return
				case t := <-tk.C:
					Try1(o, Next(t), func() { o.Error(ErrOops) })
				}
			}
		})
	}
}

// Timer creates an Observable that emits a [time.Time] value
// after a particular time span has passed, and then completes.
func Timer(d time.Duration) Observable[time.Time] {
	return func(c Context, o Observer[time.Time]) {
		tm := time.NewTimer(d)
		c.Go(func() {
			select {
			case <-c.Done():
				tm.Stop()
				o.Error(c.Cause())
			case t := <-tm.C:
				Try1(o, Next(t), func() { o.Error(ErrOops) })
				o.Complete()
			}
		})
	}
}
