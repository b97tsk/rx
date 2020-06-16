package queue_test

import (
	"testing"

	"github.com/b97tsk/rx/internal/queue"
)

func TestQueue(t *testing.T) {
	var q queue.Queue

	t.Logf("Len=%v Cap=%v", q.Len(), q.Cap())
	if q.Len() != 0 {
		t.FailNow()
	}

	for _, r := range "ABCDEFGHIJKLM" {
		r := string(r)
		q.Push(r)
		t.Logf("PushBack(%v): Len=%v Cap=%v Front=%v Back=%v", r, q.Len(), q.Cap(), q.Front(), q.Back())
		if q.Back() != r {
			t.FailNow()
		}
	}
	for _, r := range "ABCDEFG" {
		r := string(r)
		v := q.Pop()
		t.Logf("PopFront(%v): Len=%v Cap=%v Front=%v Back=%v", v, q.Len(), q.Cap(), q.Front(), q.Back())
		if v != r {
			t.FailNow()
		}
	}

	cloned := q.Clone()
	for i, r := range "HIJKLM" {
		r := string(r)
		ith := cloned.At(i)
		if ith != r {
			t.Logf("At(%v): %v expected, but got %v", i, r, ith)
			t.FailNow()
		}
	}
	cloned.Init()
	if cloned.Len() != 0 {
		t.Logf("Reset(nil): Len=%v Cap=%v", cloned.Len(), cloned.Cap())
		t.FailNow()
	}

	for _, r := range "NOPQRSTUVWXYZ" {
		r := string(r)
		q.Push(r)
		t.Logf("PushBack(%v): Len=%v Cap=%v Front=%v Back=%v", r, q.Len(), q.Cap(), q.Front(), q.Back())
		if q.Back() != r {
			t.FailNow()
		}
	}
	for _, r := range "HIJKLMNOPQRSTUVWXY" {
		r := string(r)
		v := q.Pop()
		t.Logf("PopFront(%v): Len=%v Cap=%v Front=%v Back=%v", v, q.Len(), q.Cap(), q.Front(), q.Back())
		if v != r {
			t.FailNow()
		}
	}

	v := q.Pop()
	t.Logf("PopFront(%v): Len=%v Cap=%v", v, q.Len(), q.Cap())
	if v != "Z" || q.Len() != 0 {
		t.FailNow()
	}
}
