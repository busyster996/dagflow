package task

import (
	"io"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
)

// Detail
// @Summary		详情
// @Description	获取任务详情, 支持SSE订阅
// @Tags		任务
// @Accept		application/json
// @Produce		application/json
// @Param		task path string true "任务名称"
// @Success		200 {object} base.IResponse[types.STaskRes]
// @Failure		500 {object} base.IResponse[any]
// @Router		/api/v1/task/{task} [get]
func Detail(c *gin.Context) {
	taskName := c.Param("task")
	if taskName == "" {
		base.Send(c, base.WithCode[any](base.CodeNoData).WithError(errors.New("task does not exist")))
		return
	}
	if c.GetHeader("Accept") != base.EventStreamMimeType {
		code, task, err := service.Task(taskName).Detail()
		if err != nil {
			base.Send(c, base.WithCode[any](base.CodeNoData).WithError(err))
			return
		}
		c.Request.Header.Set(types.XTaskState, task.State)
		c.Header(types.XTaskState, task.State)

		base.Send(c, base.WithData(task).WithCode(code).WithError(errors.New(task.Message)))
		return
	}
	ticker := time.NewTicker(30 * time.Second) // 每30秒发送心跳
	defer ticker.Stop()
	var lastCode base.Code
	var lastError error
	var last *types.STaskRes
	c.Stream(func(w io.Writer) bool {
		select {
		case <-ticker.C:
			c.SSEvent("heartbeat", "keepalive")
			return true
		case <-c.Done():
			return false
		default:
			code, current, err := service.Task(taskName).Detail()
			if lastCode == code && errors.Is(err, lastError) && reflect.DeepEqual(last, current) {
				time.Sleep(1 * time.Second)
				return true
			}
			c.SSEvent("message", base.WithData(current).WithError(err).WithCode(code))
			lastCode = code
			lastError = err
			last = current
			time.Sleep(1 * time.Second)
			return true
		}
	})
}
