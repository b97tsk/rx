package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/internal/atomic"
	"github.com/b97tsk/rx/internal/ctxutil"
	"github.com/b97tsk/rx/internal/queue"
)

// ReplayOptions are the options for MulticastReplay.
type ReplayOptions struct {
	BufferSize int
	WindowTime time.Duration
}

// MulticastReplay returns a Double whose Observable part takes care of all
// Observers that subscribes to it, which will receive buffered emissions and
// new emissions from Double's Observer part.
func MulticastReplay(opts *ReplayOptions) Double {
	d := &multicastReplay{}

	if opts != nil {
		d.ReplayOptions = *opts
	}

	return Double{
		Observable: d.subscribe,
		Observer:   d.sink,
	}
}

// MulticastReplayFactory returns a DoubleFactory that wraps calls to
// MulticastReplay.
func MulticastReplayFactory(opts *ReplayOptions) DoubleFactory {
	return func() Double { return MulticastReplay(opts) }
}

type multicastReplay struct {
	ReplayOptions
	multicast
	buffer     queue.Queue
	bufferRefs *atomic.Uint32s
}

type multicastReplayElement struct {
	Value    interface{}
	Deadline time.Time
}

func (d *multicastReplay) bufferForRead() (queue.Queue, *atomic.Uint32s) {
	refs := d.bufferRefs

	if refs == nil {
		refs = new(atomic.Uint32s)
		d.bufferRefs = refs
	}

	refs.Add(1)

	return d.buffer, refs
}

func (d *multicastReplay) bufferForWrite() *queue.Queue {
	refs := d.bufferRefs

	if refs != nil && !refs.Equals(0) {
		d.buffer = d.buffer.Clone()
		d.bufferRefs = nil
	}

	return &d.buffer
}

func (d *multicastReplay) trimBuffer(b *queue.Queue) {
	if d.WindowTime > 0 {
		if b == nil {
			b = d.bufferForWrite()
		}

		now := time.Now()

		for b.Len() > 0 {
			if b.Front().(multicastReplayElement).Deadline.After(now) {
				break
			}

			b.Pop()
		}
	}

	if bufferSize := d.BufferSize; bufferSize > 0 {
		if b == nil {
			b = d.bufferForWrite()
		}

		for b.Len() > bufferSize {
			b.Pop()
		}
	}
}

func (d *multicastReplay) sink(t Notification) {
	d.mu.Lock()

	switch {
	case d.err != nil:
		d.mu.Unlock()

	case t.HasValue:
		lst := d.lst.Clone()
		defer lst.Release()

		var deadline time.Time

		if windowTime := d.WindowTime; windowTime > 0 {
			deadline = time.Now().Add(windowTime)
		}

		b := d.bufferForWrite()
		b.Push(multicastReplayElement{t.Value, deadline})
		d.trimBuffer(b)

		d.mu.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}

	default:
		var lst observerList

		d.lst.Swap(&lst)

		d.err = errCompleted

		if t.HasError {
			d.err = t.Error

			if d.err == nil {
				d.err = errNil
			}

			d.buffer.Init()
		}

		d.mu.Unlock()

		for _, observer := range lst.Observers {
			observer.Sink(t)
		}
	}
}

func (d *multicastReplay) subscribe(ctx context.Context, sink Observer) {
	d.mu.Lock()

	err := d.err
	if err == nil {
		ctx, cancel := context.WithCancel(ctx)

		observer := sink.WithCancel(cancel).MutexContext(ctx)

		d.lst.Append(&observer)

		finalize := func() {
			d.mu.Lock()
			d.lst.Remove(&observer)
			d.mu.Unlock()
		}

		for d.cws == nil || !d.cws.Submit(ctx, finalize) {
			d.cws = ctxutil.NewContextWaitService()
		}
	}

	d.trimBuffer(nil)

	b, refs := d.bufferForRead()
	defer refs.Sub(1)

	d.mu.Unlock()

	for i, j := 0, b.Len(); i < j; i++ {
		if ctx.Err() != nil {
			return
		}

		sink.Next(b.At(i).(multicastReplayElement).Value)
	}

	if err != nil {
		if err == errCompleted {
			sink.Complete()

			return
		}

		if err == errNil {
			err = nil
		}

		sink.Error(err)
	}
}
