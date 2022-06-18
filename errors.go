package rx

import (
	"errors"
)

var (
	ErrEmpty     = errors.New("empty")
	ErrFinalized = errors.New("finalized")
	ErrNotSingle = errors.New("not single")
	ErrTimeout   = errors.New("timeout")
)

var (
	errCompleted = errors.New("completed")
	errNil       = errors.New("nil")
)

func errOrErrNil(err error) error {
	if err != nil {
		return err
	}

	return errNil
}

func cleanErrNil(err error) error {
	if err != errNil {
		return err
	}

	return nil
}
