package step

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/pkg/logx"
)

// Manager
// @Summary		管理
// @Description	管理指定任务的指定步骤, 支持暂停、恢复、终止、超时暂停自动恢复
// @Tags		步骤
// @Accept		application/json
// @Produce		application/json
// @Param		task path string true "任务名称"
// @Param		step path string true "步骤名称"
// @Param		action query string false "操作项" Enums(paused,kill,pause,resume) default(paused)
// @Param		duration query string false "暂停多久, 如果没设置则需要手工恢复" default(1m)
// @Success		200 {object} base.IResponse[any]
// @Failure		500 {object} base.IResponse[any]
// @Router		/api/v1/task/{task}/step/{step} [put]
func Manager(c *gin.Context) {
	taskName := c.Param("task")
	if taskName == "" {
		base.Send(c, base.WithCode[any](base.CodeNoData).WithError(errors.New("task does not exist")))
		return
	}
	stepName := c.Param("step")
	if stepName == "" {
		base.Send(c, base.WithCode[any](base.CodeNoData).WithError(errors.New("step does not exist")))
		return
	}
	action := c.DefaultQuery("action", "paused")
	duration := c.DefaultQuery("duration", "-1")
	err := service.Step(taskName, stepName).Manager(action, duration)
	if err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](base.CodeFailed).WithError(err))
		return
	}
	base.Send(c, base.WithCode[any](base.CodeSuccess))
}
