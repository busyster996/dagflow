package task

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/pkg/logx"
)

// Delete
// @Summary		删除
// @Description	删除指定任务
// @Tags		任务
// @Accept		application/json
// @Produce		application/json
// @Param		task path string true "任务名称"
// @Success		200 {object} types.SBase[any]
// @Failure		500 {object} types.SBase[any]
// @Router		/api/v1/task/{task} [delete]
func Delete(c *gin.Context) {
	taskName := c.Param("task")
	if taskName == "" {
		base.Send(c, base.WithCode[any](types.CodeNoData).WithError(errors.New("task does not exist")))
		return
	}
	err := service.Task(taskName).Delete()
	if err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](types.CodeFailed).WithError(err))
		return
	}
	base.Send(c, base.WithCode[any](types.CodeSuccess))
}
