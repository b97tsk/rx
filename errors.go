package rx

import (
	"errors"
)

var (
	ErrDropped         = errors.New("dropped")
	ErrEmpty           = errors.New("empty")
	ErrNotNotification = errors.New("not notification")
	ErrNotObservable   = errors.New("not observable")
	ErrNotSingle       = errors.New("not single")
	ErrOutOfRange      = errors.New("out of range")
	ErrTimeout         = errors.New("timeout")
)

// Errors for internal use only.
var (
	errCompleted = errors.New("completed")
	errNil       = errors.New("nil")
)
