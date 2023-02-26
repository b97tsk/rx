package waitgroup_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx/internal/waitgroup"
)

func TestWaitGroup(t *testing.T) {
	t.Parallel()

	ctx1, wg1 := waitgroup.Install(context.Background())
	ctx2, wg2 := waitgroup.Install(ctx1)

	done := make(chan struct{})

	go func() {
		wg2.Wait()
		wg1.Wait()
		close(done)
	}()

	wg1.Done()
	wg2.Done()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for done")
	}

	ctx1 = waitgroup.Hoist(ctx1)
	ctx2 = waitgroup.Hoist(ctx2)

	if wg1 != waitgroup.Get(ctx1) || wg2 != waitgroup.Get(ctx2) {
		t.Fail()
	}
}
