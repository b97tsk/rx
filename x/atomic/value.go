package atomic

import (
	"sync/atomic"
)

// A Value provides an atomic load and store of a consistently typed value.
// The zero value for a Value returns nil from Load. Once Store has been
// called, a Value must not be copied.
type Value struct {
	atomic.Value
}
