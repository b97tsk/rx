package rx

import (
	"context"
)

// A CongestingMergeConfigure is a configure for CongestingMerge.
type CongestingMergeConfigure struct {
	Project    func(interface{}, int) Observable
	Concurrent int
}

// MakeFunc creates an OperatorFunc from this type.
func (conf CongestingMergeConfigure) MakeFunc() OperatorFunc {
	return MakeFunc(congestingMergeOperator(conf).Call)
}

type congestingMergeOperator CongestingMergeConfigure

func (op congestingMergeOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	sink = Mutex(Finally(sink, cancel))

	completeSignal := make(chan struct{}, 1)

	type X struct {
		Index       int
		ActiveCount int
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			x := <-cx

			for x.ActiveCount == op.Concurrent {
				cx <- x
				select {
				case <-done:
					return
				case <-completeSignal:
				}
				x = <-cx
			}

			x.ActiveCount++

			outerIndex := x.Index
			outerValue := t.Value
			x.Index++

			// calls op.Project synchronously
			obs := op.Project(outerValue, outerIndex)

			cx <- x

			obs.Subscribe(ctx, func(t Notification) {
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

	return ctx, cancel
}

// CongestingMerge creates an output Observable which concurrently emits all
// values from every given input Observable.
//
// CongestingMerge flattens multiple Observables together by blending their
// values into one Observable.
//
// It's like Merge, but it may congest the source due to concurrent limit.
func CongestingMerge(observables ...Observable) Observable {
	return FromObservables(observables...).Pipe(operators.CongestingMergeAll())
}

// CongestingMergeAll converts a higher-order Observable into a first-order
// Observable which concurrently delivers all values that are emitted on the
// inner Observables.
//
// It's like MergeAll, but it may congest the source due to concurrent limit.
func (Operators) CongestingMergeAll() OperatorFunc {
	return operators.CongestingMergeMap(ProjectToObservable)
}

// CongestingMergeMap creates an Observable that projects each source value to
// an Observable which is merged in the output Observable.
//
// CongestingMergeMap maps each value to an Observable, then flattens all of
// these inner Observables using CongestingMergeAll.
//
// It's like MergeMap, but it may congest the source due to concurrent limit.
func (Operators) CongestingMergeMap(project func(interface{}, int) Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := congestingMergeOperator{project, -1}
		return source.Lift(op.Call)
	}
}

// CongestingMergeMapTo creates an Observable that projects each source value
// to the same Observable which is merged multiple times in the output
// Observable.
//
// It's like CongestingMergeMap, but maps each value always to the same inner
// Observable.
//
// It's like MergeMapTo, but it may congest the source due to concurrent limit.
func (Operators) CongestingMergeMapTo(inner Observable) OperatorFunc {
	return operators.CongestingMergeMap(func(interface{}, int) Observable { return inner })
}
