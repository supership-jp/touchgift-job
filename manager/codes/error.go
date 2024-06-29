package codes

import (
	"errors"
)

// ErrNoData is error when no data
var ErrNoData = errors.New("no data")

// ErrFailedUpdate is error when not updated
var ErrFailedUpdate = errors.New("failed to update")

// ErrConditionFailed is error when not updated
var ErrConditionFailed = errors.New("conditional check failed")

// ErrDoNothing　is do nothing
var ErrDoNothing = errors.New("do nothing")
