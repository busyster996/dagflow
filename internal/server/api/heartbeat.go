package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Heartbeat
// @Summary		心跳
// @Description	用于判断服务是否正常
// @Tags		默认
// @Accept		application/json
// @Produce		application/json
// @Success		200 {object} string
// @Failure		500 {object} string
// @Router		/heartbeat [get]
func Heartbeat(c *gin.Context) {
	c.AbortWithStatus(http.StatusOK)
}
