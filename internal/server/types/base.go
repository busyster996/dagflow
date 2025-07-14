package types

import (
	"encoding/json"
	"strings"

	"github.com/gorilla/websocket"
	"go.yaml.in/yaml/v3"

	"github.com/busyster996/dagflow/internal/utils"
)

const (
	XTaskName  = "X-Task-Name"
	XTaskState = "X-Task-STATE"
)

type SBase[T any] struct {
	Code      Code    `json:"code" yaml:"code"`
	Message   Message `json:"message" yaml:"message" swaggertype:"string"`
	Timestamp int64   `json:"timestamp" yaml:"timestamp"`
	Data      T       `json:"data" yaml:"data"`
}

type Message []string

func (msg Message) String() string {
	return strings.Join(utils.RemoveDuplicate(msg), "; ")
}

func (msg Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.String())
}

func (msg Message) MarshalYAML() (interface{}, error) {
	data, err := yaml.Marshal(msg.String())
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (msg *Message) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*msg = strings.Split(s, "; ")
	return nil
}

func (msg *Message) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	*msg = strings.Split(s, "; ")
	return nil
}

type Code int

const (
	CodeSuccess Code = 0
	CodeRunning Code = iota + 1000
	CodeFailed
	CodeNoData
	CodePending
	CodePaused
	CodeSkipped
)

var CodeMap = map[Code]string{
	CodeSuccess: "success",
	CodeRunning: "running",
	CodeFailed:  "failed",
	CodeNoData:  "no data",
	CodePending: "pending",
	CodePaused:  "paused",
	CodeSkipped: "skipped",
}

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
