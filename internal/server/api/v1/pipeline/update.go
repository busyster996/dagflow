package pipeline

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/pkg/logx"
)

// Update
// @Summary 	更新
// @Description 更新指定流水线
// @Tags 		流水线
// @Accept		application/json
// @Produce		application/json
// @Param		pipeline path string true "流水线名称"
// @Param		content body types.SPipelineUpdateReq true "更新内容"
// @Success		200 {object} base.IResponse[any]
// @Failure		500 {object} base.IResponse[any]
// @Router		/api/v1/pipeline/{pipeline} [post]
func Update(c *gin.Context) {
	pipelineName := c.Param("pipeline")
	if pipelineName == "" {
		base.Send(c, base.WithCode[any](base.CodeNoData).WithError(errors.New("task does not exist")))
		return
	}
	var req = new(types.SPipelineUpdateReq)
	if err := c.ShouldBind(req); err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](base.CodeFailed).WithError(err))
		return
	}
	if err := service.Pipeline(pipelineName).Update(req); err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](base.CodeFailed).WithError(err))
		return
	}
	base.Send(c, base.WithCode[any](base.CodeSuccess))
}
