package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// A CongestingMergeConfigure is a configure for CongestingMerge.
type CongestingMergeConfigure struct {
	Project    func(interface{}, int) rx.Observable
	Concurrent int
}

// Use creates an Operator from this configure.
func (configure CongestingMergeConfigure) Use() rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := congestingMergeObservable{source, configure}
		return rx.Create(obs.Subscribe)
	}
}

type congestingMergeObservable struct {
	Source rx.Observable
	CongestingMergeConfigure
}

func (obs congestingMergeObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	done := ctx.Done()
	sink = rx.Mutex(sink)

	completeSignal := make(chan struct{}, 1)

	type X struct {
		Index       int
		ActiveCount int
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			x := <-cx

			for x.ActiveCount == obs.Concurrent {
				cx <- x
				select {
				case <-done:
					return
				case <-completeSignal:
				}
				x = <-cx
			}

			x.ActiveCount++

			sourceIndex := x.Index
			sourceValue := t.Value
			x.Index++

			// calls obs.Project synchronously
			obs := obs.Project(sourceValue, sourceIndex)

			cx <- x

			obs.Subscribe(ctx, func(t rx.Notification) {
				switch {
				case t.HasValue || t.HasError:
					sink(t)
				default:
					x := <-cx
					x.ActiveCount--
					cx <- x
					select {
					case completeSignal <- struct{}{}:
					default:
					}
				}
			})

		case t.HasError:
			sink(t)

		default:
			x := <-cx
			if x.ActiveCount > 0 {
				go func() {
					for x.ActiveCount > 0 {
						cx <- x
						select {
						case <-done:
							return
						case <-completeSignal:
						}
						x = <-cx
					}
					cx <- x
					sink.Complete()
				}()
				return
			}
			cx <- x
			sink(t)
		}
	})
}

// CongestingMergeAll converts a higher-order Observable into a first-order
// Observable which concurrently delivers all values that are emitted on the
// inner Observables.
//
// It's like MergeAll, but it may congest the source due to concurrent limit.
func CongestingMergeAll() rx.Operator {
	return CongestingMergeMap(rx.ProjectToObservable)
}

// CongestingMergeMap creates an Observable that projects each source value to
// an Observable which is merged in the output Observable.
//
// CongestingMergeMap maps each value to an Observable, then flattens all of
// these inner Observables using CongestingMergeAll.
//
// It's like MergeMap, but it may congest the source due to concurrent limit.
func CongestingMergeMap(project func(interface{}, int) rx.Observable) rx.Operator {
	return CongestingMergeConfigure{project, -1}.Use()
}

// CongestingMergeMapTo creates an Observable that projects each source value
// to the same Observable which is merged multiple times in the output
// Observable.
//
// It's like CongestingMergeMap, but maps each value always to the same inner
// Observable.
//
// It's like MergeMapTo, but it may congest the source due to concurrent limit.
func CongestingMergeMapTo(inner rx.Observable) rx.Operator {
	return CongestingMergeMap(func(interface{}, int) rx.Observable { return inner })
}
