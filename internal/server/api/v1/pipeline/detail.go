package pipeline

import (
	"io"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/pkg/logx"
)

// Detail
// @Summary 	详情
// @Description 获取指定流水线详情, 支持SSE订阅
// @Tags 		流水线
// @Accept		application/json
// @Produce		application/json
// @Param		pipeline path string true "流水线名称"
// @Success		200 {object} base.IResponse[types.SPipelineRes]
// @Failure		500 {object} base.IResponse[any]
// @Router		/api/v1/pipeline/{pipeline} [get]
func Detail(c *gin.Context) {
	pipelineName := c.Param("pipeline")
	if pipelineName == "" {
		base.Send(c, base.WithCode[any](base.CodeNoData).WithError(errors.New("pipeline does not exist")))
		return
	}
	if c.GetHeader("Accept") != base.EventStreamMimeType {
		res, err := service.Pipeline(pipelineName).Detail()
		if err != nil {
			logx.Errorln(err)
			base.Send(c, base.WithCode[any](base.CodeFailed).WithError(err))
			return
		}
		base.Send(c, base.WithData(res).WithCode(base.CodeSuccess))
		return
	}
	ticker := time.NewTicker(30 * time.Second) // 每30秒发送心跳
	defer ticker.Stop()
	var last *types.SPipelineRes
	c.Stream(func(w io.Writer) bool {
		select {
		case <-ticker.C:
			c.SSEvent("heartbeat", "keepalive")
			return true
		case <-c.Done():
			return false
		default:
			current, err := service.Pipeline(pipelineName).Detail()
			if err != nil {
				logx.Errorln(err)
				c.SSEvent("error", err.Error())
				return false
			}
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
