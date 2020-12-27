package rx

import (
	"github.com/b97tsk/rx/internal/atomic"
)

type observerList struct {
	Observers []*Observer

	refs *atomic.Uint32s
}

func (lst *observerList) Clone() observerList {
	refs := lst.refs

	if refs == nil {
		refs = new(atomic.Uint32s)
		lst.refs = refs
	}

	refs.Add(1)

	return observerList{lst.Observers, refs}
}

func (lst *observerList) Release() {
	if refs := lst.refs; refs != nil {
		refs.Sub(1)

		lst.refs = nil
	}
}

func (lst *observerList) Append(observer *Observer) {
	observers := lst.Observers
	oldcap := cap(observers)

	observers = append(observers, observer)
	lst.Observers = observers

	if cap(observers) != oldcap {
		if refs := lst.refs; refs != nil && !refs.Equals(0) {
			lst.refs = nil
		}
	}
}

func (lst *observerList) Remove(observer *Observer) {
	observers := lst.Observers

	for i, sink := range observers {
		if sink == observer {
			n := len(observers)

			if refs := lst.refs; refs != nil && !refs.Equals(0) {
				newObservers := make([]*Observer, n-1, n)
				copy(newObservers, observers[:i])
				copy(newObservers[i:], observers[i+1:])
				lst.Observers = newObservers
				lst.refs = nil
			} else {
				copy(observers[i:], observers[i+1:])
				observers[n-1] = nil
				lst.Observers = observers[:n-1]
			}

			break
		}
	}
}

func (lst *observerList) Swap(other *observerList) {
	*lst, *other = *other, *lst
}
