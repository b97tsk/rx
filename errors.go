package rx

import (
	"errors"
)

var (
	ErrEmpty             = errors.New("empty")
	ErrNoOptions         = errors.New("no options")
	ErrNotNotification   = errors.New("not notification")
	ErrNotObservable     = errors.New("not observable")
	ErrNotSingle         = errors.New("not single")
	ErrOutOfRange        = errors.New("out of range")
	ErrTimeout           = errors.New("timeout")
	ErrUnsupportedOption = errors.New("unsupported option")
)
