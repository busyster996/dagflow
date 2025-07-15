package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// heartbeat
// @Summary		心跳
// @Description	用于判断服务是否正常
// @Tags		默认
// @Accept		application/json
// @Produce		application/json
// @Success		200 {object} string
// @Failure		500 {object} string
// @Router		/heartbeat [get]
func heartbeat(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}
