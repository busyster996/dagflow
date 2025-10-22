package common

import (
	"github.com/pkg/errors"
)

type ExecCode int64

const (
	ExecCodeSuccess   ExecCode = 0
	ExecCodeFailed    ExecCode = -1
	ExecCodeSkipped   ExecCode = -2
	ExecCodeKilled    ExecCode = -997
	ExecCodeTimeout   ExecCode = -998
	ExecCodeSystemErr ExecCode = -999
)

func (e ExecCode) Int64() int64 {
	return int64(e)
}

var (
	ExecErrTimeOut = errors.New("forced termination by timeout")
	ExecErrKilled  = errors.New("artificial force kill")
)

const (
	ExecConsoleStart = "OSREAPI::CONSOLE::START"
	ExecConsoleDone  = "OSREAPI::CONSOLE::DONE"
)
