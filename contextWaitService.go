package rx

import (
	"context"
	"reflect"
	"time"
)

const autoCloseDelay = 30 * time.Second

type contextWaitService chan chan<- contextWaitAction

type contextWaitAction struct {
	Context  context.Context
	Callback func()
}

func newContextWaitService() contextWaitService {
	actionChan := make(chan contextWaitAction, 1)
	service := make(contextWaitService, 1)
	service <- actionChan
	go func() {
		var (
			cases   []reflect.SelectCase
			actions []contextWaitAction
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
				<-service
				if len(actionChan) == 0 {
					close(service)
					return
				}
				service <- actionChan
				cases[0].Chan = reflect.ValueOf(nil)
				continue
			}
			if i == 1 {
				action := v.Interface().(contextWaitAction)
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
				actions[i] = contextWaitAction{}
				j := len(actions) - 1
				actions[i], actions[j] = actions[j], actions[i]
				actions = actions[:j]
				action.Callback()
			}
			if len(actions) == 0 {
				cases[0].Chan = reflect.ValueOf(time.After(autoCloseDelay))
			}
		}
	}()
	return service
}

func (service contextWaitService) Submit(ctx context.Context, cb func()) bool {
	actionChan, serviceAvailable := <-service
	if !serviceAvailable {
		return false
	}
	actionChan <- contextWaitAction{ctx, cb}
	service <- actionChan
	return true
}
