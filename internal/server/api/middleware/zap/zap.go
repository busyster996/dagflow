package zap

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/busyster996/dagflow/pkg/logx"
)

func Logger(c *gin.Context) {
	start := time.Now()
	c.Next()
	latency := time.Since(start)

	status := c.Writer.Status()
	clientIP := c.ClientIP()
	method := c.Request.Method
	proto := c.Request.Proto
	path := c.Request.URL.String()
	userAgent := c.Request.UserAgent()

	if len(c.Errors) > 0 {
		for _, err := range c.Errors.Errors() {
			logx.Errorln(clientIP, method, proto, status, path, latency, err)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		logx.Infoln(clientIP, method, proto, status, path, latency, userAgent)
	}
}

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			handlePanic(c, err)
		}
	}()
	c.Next()
}

func handlePanic(c *gin.Context, err interface{}) {
	if isBrokenPipeError(err) {
		httpRequest, _ := httputil.DumpRequest(c.Request, false)
		logx.Errorln("Broken pipe:", c.Request.URL.Path, string(httpRequest), err)
		c.Abort() // Avoid returning InternalServerError for broken pipes
		return
	}

	// Log panic details and return 500
	httpRequest, _ := httputil.DumpRequest(c.Request, false)
	logx.Errorln("[Recovery from panic]",
		time.Now().Format(time.RFC3339),
		string(httpRequest),
		string(debug.Stack()),
		err,
	)
	c.AbortWithStatus(http.StatusInternalServerError)
}

func isBrokenPipeError(err interface{}) bool {
	ne, ok := err.(*net.OpError)
	if !ok {
		return false
	}
	se, ok := ne.Err.(*os.SyscallError)
	if !ok {
		return false
	}

	errMsg := strings.ToLower(se.Error())
	return strings.Contains(errMsg, "broken pipe") || strings.Contains(errMsg, "connection reset by peer")
}
