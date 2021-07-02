package mysql

import (
	"fmt"
)

type ErrInvalidScanTo struct {
	target string
}

func (e *ErrInvalidScanTo) Error() string {
	return fmt.Sprintf(
		"invalid scan to type, need non-nil %s", e.target,
	)
}

func NewErrInvalidScanTo(target string) *ErrInvalidScanTo {
	return &ErrInvalidScanTo{
		target: target,
	}
}
