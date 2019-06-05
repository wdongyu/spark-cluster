package log

import "errors"

var (
	ErrPodPending error = errors.New("Pod is pending now")
)
