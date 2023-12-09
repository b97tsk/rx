package rx

import "errors"

var (
	ErrEmpty     = errors.New("empty")
	ErrFinalized = errors.New("finalized")
	ErrNil       = errors.New("nil")
	ErrNotSingle = errors.New("not single")
	ErrTimeout   = errors.New("timeout")
)
