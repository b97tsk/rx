package operators

import (
	"github.com/b97tsk/rx"
)

// CongestingZipAll converts a higher-order Observable into a first-order
// Observable by waiting for the outer Observable to complete, then applying
// CongestingZip.
//
// CongestingZipAll flattens an Observable-of-Observables by applying
// CongestingZip when the Observable-of-Observables completes.
//
// It's like ZipAll, but it congests subscribed Observables.
func CongestingZipAll() rx.Operator {
	return ToObservablesConfigure{rx.CongestingZip}.Use()
}
