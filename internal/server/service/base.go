package service

import (
	"fmt"
	"slices"
	"strings"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/storage/models"
)

func ConvertState(state models.State) base.Code {
	switch state {
	case models.StateStopped:
		return base.CodeSuccess
	case models.StateRunning:
		return base.CodeRunning
	case models.StatePending:
		return base.CodePending
	case models.StatePaused:
		return base.CodePaused
	case models.StateFailed:
		return base.CodeFailed
	case models.StateSkipped:
		return base.CodeSkipped
	default:
		return base.CodeNoData
	}
}

func GenerateStateMessage(baseMessage string, groups map[models.State][]string) string {
	var keys []models.State
	for k := range groups {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	var messages []string
	if baseMessage != "" {
		messages = append(messages, baseMessage)
	}
	for _, key := range keys {
		count := len(groups[key])
		messages = append(messages, fmt.Sprintf("%d %s", count, models.StateMap[key]))
	}
	return strings.Join(messages, "; ")
}
