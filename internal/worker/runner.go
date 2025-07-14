package worker

import (
	"strings"

	"github.com/busyster996/dagflow/internal/worker/runner"
)

func (s *sStep) newExecutorRunner() (runner.IRunner, error) {
	commandType, err := s.stg.Type()
	if err != nil {
		return nil, err
	}
	cmdType, subCmd, _ := strings.Cut(commandType, "@")
	if subCmd == "" {
		subCmd = cmdType
	}
	executor, err := runner.Get(cmdType)
	if err != nil {
		return nil, err
	}
	return executor(s.stg, subCmd, s.workspace, s.scriptDir)
}
