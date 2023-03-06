package rx

import (
	"errors"
)

var (
	ErrEmpty     = errors.New("empty")
	ErrFinalized = errors.New("finalized")
	ErrNil       = errors.New("nil")
	ErrNotSingle = errors.New("not single")
	ErrTimeout   = errors.New("timeout")
)

var (
	errComplete = errors.New("complete")
	errNil      = errors.New("nil")
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
