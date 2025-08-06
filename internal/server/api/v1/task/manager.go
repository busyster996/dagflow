package task

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/pkg/logx"
)

// Manager
// @Summary		管理
// @Description	管理任务, 支持暂停、恢复、终止、超时暂停自动恢复
// @Tags		任务
// @Accept		application/json
// @Produce		application/json
// @Param		task path string true "任务名称"
// @Param		action query string false "操作项" Enums(paused,kill,pause,resume) default(paused)
// @Param		duration query string false "暂停多久, 如果没设置则需要手工恢复" default(1m)
// @Success		200 {object} types.SBase[any]
// @Failure		500 {object} types.SBase[any]
// @Router		/api/v1/task/{task} [put]
func Manager(c *gin.Context) {
	taskName := c.Param("task")
	if taskName == "" {
		base.Send(c, base.WithCode[any](types.CodeNoData).WithError(errors.New("task does not exist")))
		return
	}
	action := c.DefaultQuery("action", "paused")
	duration := c.DefaultQuery("duration", "-1")
	err := service.Task(taskName).Manager(action, duration)
	if err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](types.CodeFailed).WithError(err))
		return
	}
	base.Send(c, base.WithCode[any](types.CodeSuccess))
}
