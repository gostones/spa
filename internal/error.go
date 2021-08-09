package internal

import (
	"fmt"
)

// UsageError represents incorrect usage.
//
// 0 -- success
// 1 -- general failure - any standard golang error
// 2 -- usage error defined here
type UsageError struct {
	s string
}

func (r *UsageError) Error() string {
	return r.s
}

func NewUsageError(text string) error {
	return &UsageError{
		s: text,
	}
}

func NewUsageErrorf(format string, a ...interface{}) error {
	return &UsageError{
		s: fmt.Sprintf(format, a...),
	}
}
