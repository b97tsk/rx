package ctxutil

import (
	"context"
	"reflect"
	"time"
)

const autoCloseDelay = 30 * time.Second

type ContextWaitAction struct {
	Context  context.Context
	Callback func()
}

type ContextWaitService chan chan<- ContextWaitAction

func NewContextWaitService() ContextWaitService {
	actionChan := make(chan ContextWaitAction, 1)

	service := make(ContextWaitService, 1)
	service <- actionChan

	go func() {
		var (
			cases   []reflect.SelectCase
			actions []ContextWaitAction
		)

		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(time.After(autoCloseDelay)),
		})
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(actionChan),
		})

		for {
			i, v, _ := reflect.Select(cases)
			if i == 0 {
				select {
				case <-service:
					if len(actionChan) == 0 {
						close(service)

						return
					}
					service <- actionChan
				default:
				}

				cases[0].Chan = reflect.ValueOf(nil)

				continue
			}

			if i == 1 {
				action := v.Interface().(ContextWaitAction)

				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(action.Context.Done()),
				})

				actions = append(actions, action)

				cases[0].Chan = reflect.ValueOf(nil)

				continue
			}

			{
				cases[i].Chan = reflect.ValueOf(nil)
				j := len(cases) - 1
				cases[i], cases[j] = cases[j], cases[i]
				cases = cases[:j]
			}

			{
				i := i - 2
				action := actions[i]
				actions[i] = ContextWaitAction{}
				j := len(actions) - 1
				actions[i], actions[j] = actions[j], actions[i]
				actions = actions[:j]
				go action.Callback()
			}

			if len(actions) == 0 {
				cases[0].Chan = reflect.ValueOf(time.After(autoCloseDelay))
			}
		}
	}()

	return service
}

func (service ContextWaitService) Submit(ctx context.Context, cb func()) bool {
	actionChan, serviceAvailable := <-service
	if !serviceAvailable {
		return false
	}

	actionChan <- ContextWaitAction{ctx, cb}
	service <- actionChan

	return true
}
