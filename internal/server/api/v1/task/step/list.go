package step

import (
	"context"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/pkg/logx"
)

// List
// @Summary		列表
// @Description	获取指定任务的步骤列表, 支持WS长连接
// @Tags		步骤
// @Accept		application/json
// @Produce		application/json
// @Param		task path string true "任务名称"
// @Success		200 {object} types.SBase[types.SStepsRes]
// @Failure		500 {object} types.SBase[any]
// @Router		/api/v1/task/{task}/step [get]
func List(c *gin.Context) {
	taskName := c.Param("task")
	if taskName == "" {
		base.Send(c, base.WithCode[any](types.CodeNoData).WithError(errors.New("task does not exist")))
		return
	}

	var ws *websocket.Conn
	if c.IsWebsocket() {
		var err error
		ws, err = base.Upgrade(c.Writer, c.Request)
		if err != nil {
			logx.Errorln(err)
			base.Send(c, base.WithCode[any](types.CodeNoData).WithError(err))
			return
		}
	}

	if ws == nil {
		code, list, err := service.Task(taskName).Steps()
		base.Send(c, base.WithData(list).WithError(err).WithCode(code))
		return
	}

	defer base.CloseWs(ws, "Server is shutting down")

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
	var lastCode types.Code
	var lastError error
	var lastList []*types.SStepRes
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		currentCode, currentList, currentErr := service.Task(taskName).Steps()
		if currentErr != nil && currentCode == types.CodeNoData {
			_ = ws.WriteJSON(base.WithCode[any](currentCode).WithError(currentErr))
			return
		}
		// 如果状态, 错误, 列表均没有变化, 则只发送心跳
		if lastCode == currentCode && errors.Is(currentErr, lastError) && reflect.DeepEqual(lastList, currentList) {
			err := ws.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// 保存当前数据作为上一次的数据
		lastCode = currentCode
		lastError = currentErr
		lastList = currentList

		err := ws.WriteJSON(base.WithData(currentList).WithError(currentErr).WithCode(currentCode))
		if err != nil {
			return
		}
		time.Sleep(500 * time.Millisecond)
		switch currentCode {
		case types.CodeSuccess, types.CodeFailed, types.CodeSkipped:
			return
		default:

		}
	}
}
