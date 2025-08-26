package event

import "fmt"

var (
	channel    = make(chan string, defaultBufferSize)
	observable = New[string](channel)
)

func Sendf(format string, args ...interface{}) {
	Send(fmt.Sprintf(format, args...))
}

func Send(event string) {
	select {
	case channel <- event:
	default:
		// 丢弃或处理缓冲区满时的情况
	}
}

func Subscribe() (Stream[string], int64, error) {
	return observable.Subscribe()
}

func Unsubscribe(id int64) {
	observable.Unsubscribe(id)
}
