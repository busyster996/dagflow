package libs

import (
	"reflect"

	_ "github.com/avast/retry-go/v4"
	_ "github.com/dlclark/regexp2"
	_ "github.com/dustin/go-humanize"
	_ "github.com/elastic/go-freelru"
	_ "github.com/expr-lang/expr"
	_ "github.com/fsnotify/fsnotify"
	_ "github.com/gin-contrib/cors"
	_ "github.com/gin-contrib/gzip"
	_ "github.com/gin-contrib/pprof"
	_ "github.com/gin-gonic/gin"
	_ "github.com/go-cmd/cmd"
	_ "github.com/go-git/go-git/v5"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/go-zookeeper/zk"
	_ "github.com/google/uuid"
	_ "github.com/gorilla/websocket"
	_ "github.com/imroc/req/v3"
	_ "github.com/panjf2000/ants/v2"
	_ "github.com/pelletier/go-toml/v2"
	_ "github.com/pires/go-proxyproto"
	_ "github.com/pkg/errors"
	_ "github.com/rabbitmq/amqp091-go"
	_ "github.com/redis/go-redis/v9"
	_ "github.com/robfig/cron/v3"
	_ "github.com/segmentio/ksuid"
	_ "github.com/spf13/cobra"
	_ "github.com/spf13/viper"
	_ "github.com/tidwall/gjson"
	_ "github.com/tidwall/sjson"
	_ "github.com/tklauser/ps"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"
	_ "github.com/yargevad/filepathx"
	_ "go.yaml.in/yaml/v3"
	_ "google.golang.org/grpc"

	_ "github.com/busyster996/dagflow/pkg/logx"
	_ "github.com/busyster996/dagflow/pkg/rabbitmq"
)

var Symbols = map[string]map[string]reflect.Value{}

func init() {
	for _, symbols := range []interp.Exports{
		stdlib.Symbols,
		unsafe.Symbols,
		syscall.Symbols,
		unrestricted.Symbols,
		interp.Symbols,
	} {
		for name, value := range symbols {
			Symbols[name] = value
		}
	}
}

//go:generate go install github.com/traefik/yaegi/cmd/yaegi@latest
//go:generate yaegi extract github.com/avast/retry-go/v4
//go:generate yaegi extract github.com/dlclark/regexp2
//go:generate yaegi extract github.com/dustin/go-humanize
//go:generate yaegi extract github.com/elastic/go-freelru
//go:generate yaegi extract github.com/expr-lang/expr
//go:generate yaegi extract github.com/fsnotify/fsnotify
//go:generate yaegi extract github.com/gin-contrib/cors
//go:generate yaegi extract github.com/gin-contrib/gzip
//go:generate yaegi extract github.com/gin-contrib/pprof
//go:generate yaegi extract github.com/gin-gonic/gin
//go:generate yaegi extract github.com/go-cmd/cmd
//go:generate yaegi extract github.com/go-git/go-git/v5
//go:generate yaegi extract github.com/go-sql-driver/mysql
//go:generate yaegi extract github.com/go-zookeeper/zk
//go:generate yaegi extract github.com/google/uuid
//go:generate yaegi extract github.com/gorilla/websocket
//go:generate yaegi extract github.com/imroc/req/v3
//go:generate yaegi extract github.com/panjf2000/ants/v2
//go:generate yaegi extract github.com/pelletier/go-toml/v2
//go:generate yaegi extract github.com/pires/go-proxyproto
//go:generate yaegi extract github.com/pkg/errors
//go:generate yaegi extract github.com/rabbitmq/amqp091-go
//go:generate yaegi extract github.com/redis/go-redis/v9
//go:generate yaegi extract github.com/robfig/cron/v3
//go:generate yaegi extract github.com/segmentio/ksuid
//go:generate yaegi extract github.com/spf13/cobra
//go:generate yaegi extract github.com/spf13/viper
//go:generate yaegi extract github.com/tidwall/gjson
//go:generate yaegi extract github.com/tidwall/sjson
//go:generate yaegi extract github.com/tklauser/ps
//go:generate yaegi extract github.com/busyster996/dagflow/pkg/rabbitmq
//go:generate yaegi extract github.com/busyster996/dagflow/pkg/logx
//go:generate yaegi extract github.com/yargevad/filepathx
//go:generate yaegi extract go.yaml.in/yaml/v3
//go:generate yaegi extract google.golang.org/grpc
