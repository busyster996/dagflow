package step

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
)

// Log
// @Summary		日志
// @Description	指定任务指定步骤的执行输出, 支持WS长连接
// @Tags		步骤
// @Accept		application/json
// @Produce		application/json
// @Param		task path string true "任务名称"
// @Param		step path string true "步骤名称"
// @Success		200 {object} base.IResponse[types.SStepLogsRes]
// @Failure		500 {object} base.IResponse[any]
// @Router		/api/v1/task/{task}/step/{step}/log [get]
func Log(c *gin.Context) {
	taskName := c.Param("task")
	stepName := c.Param("step")
	if taskName == "" || stepName == "" {
		base.Send(c, base.WithCode[any](base.CodeNoData).WithError(errors.New("task or step does not exist")))
		return
	}

	var ws *websocket.Conn
	if c.IsWebsocket() {
		var err error
		ws, err = base.Upgrade(c.Writer, c.Request)
		if err != nil {
			base.Send(c, base.WithCode[any](base.CodeNoData).WithError(err))
			return
		}
	}
	if ws != nil {
		var ctx, cancel = context.WithCancel(c)
		defer cancel()
		go func() {
			for {
				_, _, err := ws.ReadMessage()
				if err != nil {
					cancel()
					return
				}
			}
		}()
		defer base.CloseWs(ws, "Server is shutting down")
		err := service.Step(taskName, stepName).LogStream(ctx, func(code base.Code, data types.SStepLogsRes, err error) error {
			return ws.WriteJSON(base.WithData(data).WithCode(code).WithError(err))
		})
		if err != nil {
			_ = ws.WriteJSON(base.WithCode[any](base.CodeFailed).WithError(err))
		}
		time.Sleep(500 * time.Millisecond)
		return
	}

	code, res, err := service.Step(taskName, stepName).Log()
	base.Send(c, base.WithData(res).WithCode(code).WithError(err))
}
