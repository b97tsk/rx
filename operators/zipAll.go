package operators

import (
	"github.com/b97tsk/rx"
)

// ZipAll converts a higher-order Observable into a first-order Observable by
// waiting for the outer Observable to complete, then applying Zip.
//
// ZipAll flattens an Observable-of-Observables by applying Zip when the
// Observable-of-Observables completes.
func ZipAll() rx.Operator {
	return ToObservablesConfigure{rx.Zip}.Use()
}
