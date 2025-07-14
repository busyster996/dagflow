package service

import (
	"fmt"
	"slices"
	"strings"

	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/internal/storage/models"
)

func ConvertState(state models.State) types.Code {
	switch state {
	case models.StateStopped:
		return types.CodeSuccess
	case models.StateRunning:
		return types.CodeRunning
	case models.StatePending:
		return types.CodePending
	case models.StatePaused:
		return types.CodePaused
	case models.StateFailed:
		return types.CodeFailed
	case models.StateSkipped:
		return types.CodeSkipped
	default:
		return types.CodeNoData
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
