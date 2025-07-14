package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/gorilla/websocket"

	"github.com/busyster996/dagflow/pkg/logx"
	"github.com/busyster996/dagflow/pkg/pty"
	"github.com/busyster996/dagflow/server/types"
)

type SPtyService struct {
	ws  *websocket.Conn
	pty pty.Terminal
}

func Pty(ws *websocket.Conn) (*SPtyService, error) {
	terminal, err := pty.New("")
	if err != nil {
		logx.Errorln(err)
		return nil, err
	}
	return &SPtyService{
		pty: terminal,
		ws:  ws,
	}, nil
}

func (p *SPtyService) Run() {
	go func() {
		_, _ = io.Copy(p.pty, p)
	}()
	_, _ = io.Copy(p, p.pty)
	_, _ = p.pty.Wait(context.Background())
}

func (p *SPtyService) Read(b []byte) (n int, err error) {
	messageType, data, err := p.ws.ReadMessage()
	if err != nil {
		return
	}

	dataBuffer := bytes.Trim(data, "\x00")
	_, ok := types.WebsocketMessageType[messageType]
	if !ok {
		return
	}
	if messageType == websocket.BinaryMessage {
		if dataBuffer[0] == 1 {
			ttySize := &types.STTYSize{}
			resizeMessage := bytes.Trim(dataBuffer[1:], " \n\r\t\x00\x01")
			if err = json.Unmarshal(resizeMessage, ttySize); err != nil {
				return
			}
			_ = p.pty.Resize(int16(ttySize.Rows), int16(ttySize.Cols))
			return
		}
	}

	return copy(b, dataBuffer), nil
}

func (p *SPtyService) Write(b []byte) (n int, err error) {
	n = len(b)
	err = p.ws.WriteMessage(websocket.BinaryMessage, b)
	return
}
