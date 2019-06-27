package rx

import (
	"github.com/b97tsk/rx/x/atomic"
)

type observerList struct {
	observers []*Observer
	refs      *atomic.Uint32
}

func (list *observerList) AddRef() ([]*Observer, func()) {
	refs := list.refs
	if refs == nil {
		refs = new(atomic.Uint32)
		list.refs = refs
	}
	refs.Add(1)
	return list.observers, func() { refs.Sub(1) }
}

func (list *observerList) Append(observers ...*Observer) {
	if len(observers) == 0 {
		return
	}
	if list.refs == nil || list.refs.Equals(0) {
		list.observers = append(list.observers, observers...)
		return
	}
	n := len(list.observers)
	list.observers = append(list.observers[:n:n], observers...)
	list.refs = nil
}

func (list *observerList) Remove(observer *Observer) {
	for i, sink := range list.observers {
		if sink == observer {
			observers := list.observers
			if list.refs != nil && !list.refs.Equals(0) {
				observers = append([]*Observer(nil), observers...)
				list.refs = nil
			}
			copy(observers[i:], observers[i+1:])
			n := len(observers)
			observers[n-1] = nil
			list.observers = observers[:n-1]
			break
		}
	}
}

func (list *observerList) Swap(observers []*Observer) []*Observer {
	observers, list.observers = list.observers, observers
	if list.refs != nil && !list.refs.Equals(0) {
		list.refs = nil
	}
	return observers
}