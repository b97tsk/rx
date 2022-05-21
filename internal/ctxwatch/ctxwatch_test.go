package ctxwatch_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/ctxwatch"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	v := atomic.FromInt32(0)

	cancels := []context.CancelFunc(nil)
	done := func(context.Context) { v.Sub(1) }
	add := func(n int) {
		v.Add(int32(n))

		for i := 0; i < n; i++ {
			ctx, cancel := context.WithCancel(context.Background())
			cancels = append(cancels, cancel)

			ctxwatch.Add(ctx, done)

			time.Sleep(30 * time.Millisecond)
		}
	}

	remove := func(n int) {
		for _, cancel := range cancels[:n] {
			cancel()
		}

		cancels = cancels[n:]

		time.Sleep(30 * time.Millisecond * time.Duration(n))
	}

	for n, m := 10, 50; n < m; n += 10 {
		add(m - n)
		remove(n)
	}

	if !v.Equal(0) {
		t.Fail()
	}
}
