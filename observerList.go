package rx

import (
	"github.com/b97tsk/rx/internal/atomic"
)

type observerList struct {
	Observers []*Observer

	refs *atomic.Uint32
}

func (lst *observerList) Clone() observerList {
	refs := lst.refs
	if refs == nil {
		refs = new(atomic.Uint32)
		lst.refs = refs
	}
	refs.Add(1)
	return observerList{lst.Observers, refs}
}

func (lst *observerList) Release() {
	refs := lst.refs
	if refs != nil {
		refs.Sub(1)
		lst.refs = nil
	}
}

func (lst *observerList) Append(observer *Observer) {
	refs := lst.refs
	if refs == nil || refs.Equals(0) {
		lst.Observers = append(lst.Observers, observer)
		return
	}
	observers := lst.Observers
	n := len(observers)
	lst.Observers = append(observers[:n:n], observer)
	lst.refs = nil
}

func (lst *observerList) Remove(observer *Observer) {
	observers := lst.Observers
	for i, sink := range observers {
		if sink == observer {
			n := len(observers)
			refs := lst.refs
			if refs == nil || refs.Equals(0) {
				copy(observers[i:], observers[i+1:])
				observers[n-1] = nil
				lst.Observers = observers[:n-1]
			} else {
				newObservers := make([]*Observer, n-1, n)
				copy(newObservers, observers[:i])
				copy(newObservers[i:], observers[i+1:])
				lst.Observers = newObservers
				lst.refs = nil
			}
			break
		}
	}
}

func (lst *observerList) Swap(other *observerList) {
	*lst, *other = *other, *lst
}
