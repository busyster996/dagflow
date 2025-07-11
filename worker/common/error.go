package common

import (
	"github.com/pkg/errors"
)

var (
	ErrTimeOut = errors.New("forced termination by timeout")
	ErrManual  = errors.New("artificial force termination")
)
