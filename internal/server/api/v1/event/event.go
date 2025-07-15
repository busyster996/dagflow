package event

import (
	"context"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
)

// Stream
// @Summary 	事件
// @Description 订阅系统事件, 仅支持SSE订阅
// @Tags 		事件
// @Accept		application/json
// @Produce		application/json
// @Success		200 {object} types.SBase[any]
// @Failure		500 {object} types.SBase[any]
// @Router		/api/v1/event [get]
func Stream(c *gin.Context) {
	// 判断是否为SSE连接
	if c.GetHeader("Accept") != base.EventStreamMimeType {
		base.Send(c, base.WithCode[any](types.CodeFailed).WithError(io.EOF))
		return
	}
	ctx, cancel := context.WithCancel(c)
	defer cancel()
	var event = make(chan string, 65534)
	err := service.Event().Subscribe(ctx, event)
	if err != nil {
		base.Send(c, base.WithCode[any](types.CodeFailed).WithError(err))
		return
	}
	ticker := time.NewTicker(30 * time.Second) // 每30秒发送心跳
	defer ticker.Stop()
	c.Stream(func(w io.Writer) bool {
		select {
		case e, ok := <-event:
			if !ok {
				return false
			}
			c.SSEvent("message", e)
			return true
		case <-ticker.C:
			c.SSEvent("heartbeat", "keepalive")
			return true
		case <-c.Done():
			return false
		}
	})
}
