package pool

import (
	"github.com/busyster996/dagflow/pkg/logx"
	"github.com/gin-gonic/gin"

	"github.com/busyster996/dagflow/server/api/base"
	"github.com/busyster996/dagflow/server/service"
	"github.com/busyster996/dagflow/server/types"
)

// Post
// @Summary		设置
// @Description	设置工作池大小
// @Tags		工作池
// @Accept		application/json
// @Produce		application/json
// @Param		setting body types.SPoolReq true "pool setting"
// @Success		200 {object} types.SBase[types.SPoolReq]
// @Failure		500 {object} types.SBase[any]
// @Router		/api/v1/pool [post]
func Post(c *gin.Context) {
	var req = new(types.SPoolReq)
	if err := c.ShouldBind(&req); err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](types.CodeFailed).WithError(err))
		return
	}
	res, err := service.Pool().Set(req.Size)
	if err != nil {
		base.Send(c, base.WithCode[any](types.CodeFailed).WithError(err))
		return
	}
	base.Send(c, base.WithData(res))
}
