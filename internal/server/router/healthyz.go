package router

import (
	"github.com/gin-gonic/gin"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/types"
)

// healthyz
// @Summary		健康
// @Description 用于检测服务是否正常
// @Tags		默认
// @Accept		application/json
// @Produce		application/json
// @Success		200 {object} types.SBase[types.SHealthyz]
// @Failure		500 {object} types.SBase[any]
// @Router		/healthyz [get]
func healthyz(c *gin.Context) {
	base.Send(c, base.WithData(&types.SHealthyz{
		Server: c.Request.Host,
		Client: c.Request.RemoteAddr,
		State:  "Running",
	}))
}
