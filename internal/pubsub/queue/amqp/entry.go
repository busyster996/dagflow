package amqp

import (
	"fmt"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/busyster996/dagflow/internal/pubsub/queue"
	"github.com/busyster996/dagflow/pkg/logx"
	"github.com/busyster996/dagflow/pkg/rabbitmq"
)

func init() {
	queue.Register("amqp", func(rawURL string) (queue.IBroker, error) {
		a := &sAmqp{
			publisherMap: make(map[string]*rabbitmq.Publisher),
			consumerMap:  make(map[string]*rabbitmq.Consumer),
		}
		hostname, err := os.Hostname()
		if err != nil {
			hostname = fmt.Sprintf("unknown-%d", time.Now().Nanosecond())
		}
		table := amqp091.NewConnectionProperties()
		table["connection_name"] = hostname
		a.conn, err = rabbitmq.NewConn(
			rawURL,
			rabbitmq.WithConnectionOptionsLogger(logx.GetSubLoggerWithOption(zap.AddCallerSkip(-1))),
			rabbitmq.WithConnectionOptionsConfig(rabbitmq.Config{
				Properties: table,
			}),
		)
		if err != nil {
			return nil, err
		}

		return a, nil
	})
}
