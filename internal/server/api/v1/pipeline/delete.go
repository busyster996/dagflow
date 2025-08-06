package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/pkg/logx"
)

// Delete
// @Summary 	删除
// @Description 删除指定流水线
// @Tags 		流水线
// @Accept		application/json
// @Produce		application/json
// @Param		pipeline path string true "流水线名称"
// @Success		200 {object} types.SBase[any]
// @Failure		500 {object} types.SBase[any]
// @Router		/api/v1/pipeline/{pipeline} [delete]
func Delete(c *gin.Context) {
	pipelineName := c.Param("pipeline")
	if pipelineName == "" {
		base.Send(c, base.WithCode[any](types.CodeNoData).WithError(errors.New("pipeline does not exist")))
		return
	}
	err := service.Pipeline(pipelineName).Delete()
	if err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](types.CodeFailed).WithError(err))
		return
	}
	base.Send(c, base.WithCode[any](types.CodeSuccess))
}
