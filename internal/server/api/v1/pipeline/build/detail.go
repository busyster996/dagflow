package build

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
// @Description 获取流水线指定构建任务详情, 支持SSE订阅
// @Tags 		构建
// @Accept		application/json
// @Produce		application/json
// @Param		pipeline path string true "流水线名称"
// @Param		build path string true "构建名称"
// @Success		200 {object} types.SBase[types.SPipelineBuildRes]
// @Failure		500 {object} types.SBase[any]
// @Router		/api/v1/pipeline/{pipeline}/build/{build} [get]
func Detail(c *gin.Context) {
	pipelineName := c.Param("pipeline")
	if pipelineName == "" {
		base.Send(c, base.WithCode[any](types.CodeNoData).WithError(errors.New("pipeline does not exist")))
		return
	}
	buildName := c.Param("build")
	if buildName == "" {
		base.Send(c, base.WithCode[any](types.CodeNoData).WithError(errors.New("build does not exist")))
		return
	}
	if c.GetHeader("Accept") != base.EventStreamMimeType {
		code, build, err := service.Pipeline(pipelineName).BuildDetail(buildName)
		if err != nil {
			logx.Errorln(err)
			base.Send(c, base.WithCode[any](types.CodeFailed).WithError(err))
			return
		}
		base.Send(c, base.WithData(build).WithCode(code))
		return
	}
	ticker := time.NewTicker(30 * time.Second) // 每30秒发送心跳
	defer ticker.Stop()

	var lastCode types.Code
	var lastError error
	var last *types.SPipelineBuildRes
	c.Stream(func(w io.Writer) bool {
		select {
		case <-ticker.C:
			c.SSEvent("heartbeat", "keepalive")
			return true
		case <-c.Done():
			return false
		default:
			code, current, err := service.Pipeline(pipelineName).BuildDetail(buildName)
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
