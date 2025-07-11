package pool

import (
	"io"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/busyster996/dagflow/server/api/base"
	"github.com/busyster996/dagflow/server/service"
	"github.com/busyster996/dagflow/server/types"
)

// Detail
// @Summary		详情, 支持SSE订阅
// @Description	获取工作池信息
// @Tags		工作池
// @Accept		application/json
// @Produce		application/json
// @Success		200 {object} types.SBase[types.SPoolRes]
// @Failure		500 {object} types.SBase[any]
// @Router		/api/v1/pool [get]
func Detail(c *gin.Context) {
	if c.GetHeader("Accept") != base.EventStreamMimeType {
		base.Send(c, base.WithData(service.Pool().Get()))
		return
	}

	ticker := time.NewTicker(30 * time.Second) // 每30秒发送心跳
	defer ticker.Stop()

	var last *types.SPoolRes
	c.Stream(func(w io.Writer) bool {
		select {
		case <-ticker.C:
			c.SSEvent("heartbeat", "keepalive")
			return true
		case <-c.Done():
			return false
		default:
			current := service.Pool().Get()
			if reflect.DeepEqual(last, current) {
				time.Sleep(1 * time.Second)
				return true
			}
			c.SSEvent("message", base.WithData(current))
			last = current
			time.Sleep(1 * time.Second)
			return true
		}
	})
}
