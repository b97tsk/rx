package rx

import "errors"

var (
	ErrBufferOverflow = errors.New("buffer overflow")
	ErrEmpty          = errors.New("empty")
	ErrFinalized      = errors.New("finalized")
	ErrNil            = errors.New("nil")
	ErrNotSingle      = errors.New("not single")
	ErrOops           = errors.New("oops")
	ErrTimeout        = errors.New("timeout")
	ErrUnicast        = errors.New("unicast")
)
