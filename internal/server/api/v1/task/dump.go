package task

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.yaml.in/yaml/v3"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
)

// Dump
// @Summary		导出
// @Description	导出任务
// @Tags		任务
// @Accept		application/json
// @Produce		application/json
// @Param		task path string true "任务名称"
// @Success		200 {object} types.SBase[any]
// @Failure		500 {object} types.SBase[any]
// @Router		/api/v1/task/{task}/dump [get]

func Dump(c *gin.Context) {
	taskName := c.Param("task")
	if taskName == "" {
		base.Send(c, base.WithCode[any](types.CodeNoData).WithError(errors.New("task does not exist")))
		return
	}
	res, err := service.Task(taskName).Dump()
	if err != nil {
		base.Send(c, base.WithCode[any](types.CodeFailed).WithError(err))
		return
	}
	data, err := yaml.Marshal(res)
	if err != nil {
		base.Send(c, base.WithCode[any](types.CodeFailed).WithError(err))
		return
	}
	base.Send(c, base.WithData[string](string(data)))
}
