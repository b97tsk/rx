package rx

import (
	"errors"
)

var (
	ErrEmpty           = errors.New("empty")
	ErrNotNotification = errors.New("not notification")
	ErrNotObservable   = errors.New("not observable")
	ErrNotSingle       = errors.New("not single")
	ErrOutOfRange      = errors.New("out of range")
	ErrTimeout         = errors.New("timeout")
)

var errCompleted = errors.New("completed") // Internal use only.
