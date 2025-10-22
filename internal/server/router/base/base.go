package base

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/busyster996/dagflow/internal/utility"
)

const (
	EventStreamMimeType = "text/event-stream"
	MessageSeparator    = "; " // 提取分隔符为常量
)

// IResponse 更清晰的接口命名
type IResponse[T any] interface {
	WithCode(code Code) IResponse[T]
	WithError(err error) IResponse[T]
	WithData(data T) IResponse[T]
	WithMessage(msg string) IResponse[T]
	GetCode() Code
	GetData() T
}

// Send 发送响应，支持内容协商
func Send(g *gin.Context, v interface{}) {
	contentType := g.NegotiateFormat(binding.MIMEJSON, binding.MIMEYAML, binding.MIMEYAML2)
	switch contentType {
	case binding.MIMEJSON:
		g.JSON(http.StatusOK, v)
	case binding.MIMEYAML, binding.MIMEYAML2:
		g.YAML(http.StatusOK, v)
	default:
		g.JSON(http.StatusOK, v)
	}
}

// Messages 消息列表类型
type Messages []string

// String 将消息列表转为字符串
func (msg Messages) String() string {
	if len(msg) == 0 {
		return ""
	}
	return strings.Join(utility.RemoveDuplicate(msg), MessageSeparator)
}

// MarshalJSON 自定义 JSON 序列化
func (msg Messages) MarshalJSON() ([]byte, error) {
	return json.Marshal(msg.String())
}

// MarshalYAML 自定义 YAML 序列化
func (msg Messages) MarshalYAML() (interface{}, error) {
	return msg.String(), nil
}

// UnmarshalJSON 自定义 JSON 反序列化
func (msg *Messages) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// 处理空字符串
	if s == "" {
		*msg = Messages{}
		return nil
	}

	*msg = strings.Split(s, MessageSeparator)
	return nil
}

// UnmarshalYAML 自定义 YAML 反序列化
func (msg *Messages) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	// 处理空字符串
	if s == "" {
		*msg = Messages{}
		return nil
	}

	*msg = strings.Split(s, MessageSeparator)
	return nil
}

// Code 响应状态码
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

// codeMessageMap 状态码对应的消息映射
var codeMessageMap = map[Code]string{
	CodeSuccess: "success",
	CodeRunning: "running",
	CodeFailed:  "failed",
	CodeNoData:  "no data",
	CodePending: "pending",
	CodePaused:  "paused",
	CodeSkipped: "skipped",
}

// Response 响应结构体（更清晰的命名）
type Response[T any] struct {
	Code      Code     `json:"code" yaml:"code"`
	Message   Messages `json:"message" yaml:"message" swaggertype:"string"`
	Timestamp int64    `json:"timestamp" yaml:"timestamp"`
	Data      T        `json:"data" yaml:"data"`
}

// getCodeMessage 根据状态码获取默认消息
func (r *Response[T]) getCodeMessage() string {
	if msg, ok := codeMessageMap[r.Code]; ok {
		return msg
	}
	return codeMessageMap[CodeNoData]
}

// updateTimestamp 更新时间戳
func (r *Response[T]) updateTimestamp() {
	r.Timestamp = time.Now().UnixNano()
}

// isCodeMessage 检查消息是否为状态码默认消息
func isCodeMessage(message string) bool {
	for _, msg := range codeMessageMap {
		if msg == message {
			return true
		}
	}
	return false
}

// filterCodeMessages 过滤掉状态码默认消息，保留自定义消息
func (r *Response[T]) filterCodeMessages() Messages {
	if len(r.Message) == 0 {
		return Messages{}
	}

	filtered := make(Messages, 0, len(r.Message))
	for _, msg := range r.Message {
		if !isCodeMessage(msg) {
			filtered = append(filtered, msg)
		}
	}
	return filtered
}

// WithCode 设置状态码
func (r *Response[T]) WithCode(code Code) IResponse[T] {
	// 兼容 HTTP 状态码
	if code == http.StatusOK {
		code = CodeSuccess
	}

	r.Code = code
	r.Message = append(r.filterCodeMessages(), r.getCodeMessage())
	r.updateTimestamp()
	return r
}

// WithError 设置错误信息
func (r *Response[T]) WithError(err error) IResponse[T] {
	if err == nil {
		return r.WithCode(CodeSuccess)
	}

	// 确保状态码为失败
	if r.Code == CodeSuccess {
		r.Code = CodeFailed
	}

	errMsg := strings.TrimSpace(err.Error())
	if errMsg != "" {
		r.Message = append(r.filterCodeMessages(), errMsg)
	}
	r.updateTimestamp()
	return r
}

// WithData 设置响应数据
func (r *Response[T]) WithData(data T) IResponse[T] {
	// 如果消息为空，添加默认消息
	if len(r.Message) == 0 {
		r.Message = Messages{r.getCodeMessage()}
	}

	r.Data = data
	r.updateTimestamp()
	return r
}

// WithMessage 添加自定义消息
func (r *Response[T]) WithMessage(msg string) IResponse[T] {
	msg = strings.TrimSpace(msg)
	if msg != "" {
		r.Message = append(r.filterCodeMessages(), msg)
	}
	r.updateTimestamp()
	return r
}

// GetCode 获取状态码
func (r *Response[T]) GetCode() Code {
	return r.Code
}

// GetData 获取数据
func (r *Response[T]) GetData() T {
	return r.Data
}

// NewResponse 创建新的响应对象
func NewResponse[T any]() IResponse[T] {
	return &Response[T]{
		Code:      CodeSuccess,
		Message:   Messages{codeMessageMap[CodeSuccess]},
		Timestamp: time.Now().UnixNano(),
	}
}

// WithCode 创建带状态码的响应
func WithCode[T any](code Code) IResponse[T] {
	return NewResponse[T]().WithCode(code)
}

// WithData 创建带数据的成功响应
func WithData[T any](data T) IResponse[T] {
	return NewResponse[T]().WithData(data)
}

// WithError 创建错误响应
func WithError[T any](err error) IResponse[T] {
	r := &Response[T]{
		Code:      CodeFailed,
		Timestamp: time.Now().UnixNano(),
	}
	return r.WithError(err)
}

// Success 创建成功响应的便捷方法
func Success[T any](data T) IResponse[T] {
	return WithData(data)
}

// Fail 创建失败响应的便捷方法
func Fail[T any](message string) IResponse[T] {
	return NewResponse[T]().WithCode(CodeFailed).WithMessage(message)
}
