package queue_test

import (
	"testing"

	"github.com/b97tsk/rx/x/queue"
)

func TestQueue(t *testing.T) {
	var q queue.Queue

	t.Logf("Len=%v Cap=%v", q.Len(), q.Cap())
	if q.Len() != 0 {
		t.FailNow()
	}

	q.Init()
	if q.Len() != 0 {
		t.Logf("Init: Len=%v Cap=%v", q.Len(), q.Cap())
		t.FailNow()
	}

	for _, r := range "ABCDEFGHIJKLM" {
		r := string(r)
		q.PushBack(r)
		t.Logf("PushBack(%v): Len=%v Cap=%v Front=%v Back=%v", r, q.Len(), q.Cap(), q.Front(), q.Back())
		if q.Back() != r {
			t.FailNow()
		}
	}
	for _, r := range "ABCDEFG" {
		r := string(r)
		popped := q.PopFront()
		t.Logf("PopFront(%v): Len=%v Cap=%v Front=%v Back=%v", popped, q.Len(), q.Cap(), q.Front(), q.Back())
		if popped != r {
			t.FailNow()
		}
	}
	for i, r := range "HIJKLM" {
		r := string(r)
		ith := q.At(i)
		if ith != r {
			t.Logf("At(%v): %v expected, but got %v", i, r, ith)
			t.FailNow()
		}
	}
	for _, r := range "NOPQRSTUVWXYZ" {
		r := string(r)
		q.PushBack(r)
		t.Logf("PushBack(%v): Len=%v Cap=%v Front=%v Back=%v", r, q.Len(), q.Cap(), q.Front(), q.Back())
		if q.Back() != r {
			t.FailNow()
		}
	}
	for _, r := range "HIJKLMNOPQRSTUVWXY" {
		r := string(r)
		popped := q.PopFront()
		t.Logf("PopFront(%v): Len=%v Cap=%v Front=%v Back=%v", popped, q.Len(), q.Cap(), q.Front(), q.Back())
		if popped != r {
			t.FailNow()
		}
	}

	popped := q.PopFront()
	t.Logf("PopFront(%v): Len=%v Cap=%v", popped, q.Len(), q.Cap())
	if popped != "Z" || q.Len() != 0 {
		t.FailNow()
	}
}
