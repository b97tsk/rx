package operators

import (
	"time"

	"github.com/b97tsk/rx"
)

func share(source rx.Observable) rx.Observable {
	return source.Multicast(rx.NewSubject).RefCount()
}

// Share returns a new Observable that multicasts (shares) the original
// Observable. When subscribed multiple times, it guarantees that only one
// subscription is made to the source Observable at the same time. When all
// subscribers have unsubscribed it will unsubscribe from the source Observable.
func Share() rx.Operator {
	return share
}

// ShareReplay is like Share, but it uses a ReplaySubject instead.
func ShareReplay(bufferSize int, windowTime time.Duration) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		newSubject := func() rx.Subject {
			return rx.NewReplaySubject(bufferSize, windowTime).Subject
		}
		return source.Multicast(newSubject).RefCount()
	}
}
