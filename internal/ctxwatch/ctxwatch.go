package ctxwatch

import (
	"context"
	"reflect"
	"sync"

	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/queue"
)

type watchItem struct {
	Context  context.Context
	Callback func()
}

type watchService chan<- *watchItem

func startService() watchService {
	cin, cout := make(chan *watchItem), make(chan *watchItem)

	go func() {
		var queue queue.Queue

		for {
			var (
				in   <-chan *watchItem = cin
				out  chan<- *watchItem
				outv *watchItem
			)

			if queue.Len() > 0 {
				out, outv = cout, queue.Front().(*watchItem)
			}

			select {
			case item := <-in:
				queue.Push(item)
			case out <- outv:
				queue.Pop()
			}
		}
	}()

	go func() {
		cases := []reflect.SelectCase{{}}
		itemCounter := atomic.FromInt32(0)
		workloadPerWorker := atomic.FromInt32(5)
		oldWorkloadPerWorker := 3

		for item := range cout {
			itemValue := reflect.ValueOf(item)

			for i, j := 1, len(cases); i < j; i++ {
				cases[i].Send = itemValue
			}

			dir0 := reflect.SelectSend
			if n := len(cases); int(itemCounter.Load()) >= (n-1)*oldWorkloadPerWorker {
				dir0 = reflect.SelectDefault
			}

			cases[0].Dir = dir0

			itemCounter.Add(1)

			for {
				if i, _, _ := reflect.Select(cases); i > 0 {
					break
				}

				cases[0].Dir = reflect.SelectSend

				worker := make(chan *watchItem)
				workerValue := reflect.ValueOf(worker)

				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectSend,
					Chan: workerValue,
					Send: itemValue,
				})

				w := workloadPerWorker.Load()
				workloadPerWorker.Store(w + int32(oldWorkloadPerWorker))
				oldWorkloadPerWorker = int(w)

				go startWorker(workerValue, &itemCounter, &workloadPerWorker)
			}
		}
	}()

	return cin
}

func startWorker(workerChan reflect.Value, itemCounter *atomic.Int32, workloadPerWorker *atomic.Int32) {
	cases := []reflect.SelectCase{{Dir: reflect.SelectRecv}}
	items := []*watchItem{nil}

	for {
		var chan0 reflect.Value
		if len(items)-1 < int(workloadPerWorker.Load()) {
			chan0 = workerChan
		}

		cases[0].Chan = chan0

		switch i, v, _ := reflect.Select(cases); i {
		case 0:
			item := v.Interface().(*watchItem)

			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(item.Context.Done()),
			})

			items = append(items, item)
		default:
			j := len(cases) - 1
			item := items[i]

			cases[i].Chan = cases[j].Chan
			cases[j].Chan = reflect.Value{}
			cases = cases[:j]

			items[i] = items[j]
			items[j] = nil
			items = items[:j]

			go item.Callback()

			itemCounter.Sub(1)
		}
	}
}

var shared struct {
	sync.Once
	service watchService
}

func sharedInit() {
	shared.service = startService()
}

func getService() watchService {
	shared.Do(sharedInit)
	return shared.service
}

// Add starts a background service that it waits until ctx has been cancelled,
// then it calls f in a goroutine. Successive calls use the same service to
// deal with arbitrary number of context.Contexts.
func Add(ctx context.Context, f func()) {
	getService() <- &watchItem{ctx, f}
}
