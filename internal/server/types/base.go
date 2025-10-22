package types

import (
	"github.com/gorilla/websocket"
)

const (
	XTaskName  = "X-Task-Name"
	XTaskState = "X-Task-STATE"
)

var WebsocketMessageType = map[int]string{
	websocket.BinaryMessage: "binary",
	websocket.TextMessage:   "text",
	websocket.CloseMessage:  "close",
	websocket.PingMessage:   "ping",
	websocket.PongMessage:   "pong",
}

type STTYSize struct {
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
	X    uint16 `json:"x"`
	Y    uint16 `json:"y"`
}

type SPageRes struct {
	Current int64 `json:"current" yaml:"current"`
	Size    int64 `json:"size" yaml:"size"`
	Total   int64 `json:"total" yaml:"total"`
}

type SPageReq struct {
	Page   int64  `json:"page" query:"page" yaml:"page"`
	Size   int64  `json:"size" query:"size" yaml:"size"`
	Prefix string `json:"prefix" query:"prefix" yaml:"prefix"`
}

type STimeRes struct {
	Start string `json:"start,omitempty" yaml:"start,omitempty"`
	End   string `json:"end,omitempty" yaml:"end,omitempty"`
}

type SEnvs []*SEnv

type SEnv struct {
	Name  string `json:"name" yaml:"name" binding:"required"`
	Value string `json:"value" yaml:"value"`
}
