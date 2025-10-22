package task

import (
	"github.com/gin-gonic/gin"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/pkg/logx"
)

// Post
// @Summary		创建
// @Description	创建任务
// @Tags		任务
// @Accept		application/json
// @Produce		application/json
// @Param		task body types.STaskReq true "任务内容"
// @Success		200 {object} base.IResponse[types.STaskCreateRes]
// @Failure		500 {object} base.IResponse[any]
// @Router		/api/v1/task [post]
func Post(c *gin.Context) {
	var req = new(types.STaskReq)
	if err := c.ShouldBind(req); err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](base.CodeFailed).WithError(err))
		return
	}

	if err := service.Task(req.Name).Create(req); err != nil {
		base.Send(c, base.WithCode[any](base.CodeFailed).WithError(err))
		return
	}

	c.Request.Header.Set(types.XTaskName, req.Name)
	c.Header(types.XTaskName, req.Name)

	base.Send(c, base.WithData(&types.STaskCreateRes{
		Name:  req.Name,
		Count: int64(len(req.Step)),
	}))
}
