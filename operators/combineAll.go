package operators

import (
	"github.com/b97tsk/rx"
)

// CombineAll converts a higher-order Observable into a first-order Observable
// by waiting for the outer Observable to complete, then applying CombineLatest.
//
// CombineAll flattens an Observable-of-Observables by applying CombineLatest
// when the Observable-of-Observables completes.
func CombineAll() rx.Operator {
	return ToObservablesConfigure{rx.CombineLatest}.Use()
}
