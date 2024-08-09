package queue_test

import (
	"testing"

	"github.com/b97tsk/rx/internal/queue"
)

func TestQueue(t *testing.T) {
	t.Parallel()

	var q queue.Queue[string]

	q.Init()

	t.Logf("Len=%v Cap=%v", q.Len(), q.Cap())

	if q.Len() != 0 {
		t.FailNow()
	}

	pushLetters := func(s string) {
		for _, r := range s {
			r := string(r)
			q.Push(r)

			t.Logf("Push(%v): Len=%v Cap=%v Front=%v Back=%v", r, q.Len(), q.Cap(), q.Front(), q.Back())

			if q.Back() != r {
				t.FailNow()
			}
		}
	}

	popLetters := func(s string) {
		for _, r := range s {
			r := string(r)
			v := q.Pop()

			t.Logf("Pop(%v): Len=%v Cap=%v Front=%v Back=%v", v, q.Len(), q.Cap(), q.Front(), q.Back())

			if v != r {
				t.FailNow()
			}
		}
	}

	pushLetters("ABCDEF")
	popLetters("ABCD")
	pushLetters("GHIJKL")

	for i, r := range "EFGHIJKL" {
		r := string(r)
		ith := q.At(i)

		if ith != r {
			t.Logf("At(%v): %v expected, but got %v", i, r, ith)
			t.FailNow()
		}
	}

	popLetters("EFGH")
	pushLetters("MNOPQRSTUVWXYZ")
	popLetters("IJKLMNOPQRSTUVWXY")

	cloned := q.Clone()

	if cloned.Len() != 1 || cloned.Front() != "Z" {
		t.FailNow()
	}

	cloned.Clear()

	if cloned.Len() != 0 {
		t.Logf("Clear(): Len=%v Cap=%v", cloned.Len(), cloned.Cap())
		t.FailNow()
	}

	v := q.Pop()

	t.Logf("Pop(%v): Len=%v Cap=%v", v, q.Len(), q.Cap())

	if v != "Z" || q.Len() != 0 {
		t.FailNow()
	}

	shouldPanic(t, func() { _ = q.Pop() }, "Pop called on empty Queue")
	shouldPanic(t, func() { _ = q.Front() }, "Front called on empty Queue")
	shouldPanic(t, func() { _ = q.Back() }, "Back called on empty Queue")
	shouldPanic(t, func() { _ = q.At(0) }, "At called with index out of range")
}

func shouldPanic(t *testing.T, f func(), name string) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Log(name, "did not panic.")
			t.Fail()
		}
	}()

	f()
}
