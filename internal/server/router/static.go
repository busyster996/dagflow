package router

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

//go:embed static
var staticFS embed.FS
var staticCache = sync.Map{}

func staticHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "" || path == "/" {
			path = "index.html"
		} else {
			path = strings.TrimPrefix(path, "/")
		}
		val, ok := staticCache.Load(path)
		if !ok {
			var err error
			val, err = staticFileContent(path)
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			staticCache.Store(path, val)
		}
		content := val.([]byte)
		c.Header("Content-Length", strconv.Itoa(len(content)))
		mimeType := mime.TypeByExtension(filepath.Ext(path))
		if mimeType != "" {
			c.Header("Content-Type", mimeType)
		}
		c.Status(200)
		_, _ = c.Writer.Write(content)
	}
}

var (
	once    sync.Once
	fileSys fs.FS
	initErr error
)

func init() {
	once.Do(func() {
		fileSys, initErr = fs.Sub(staticFS, "static")
	})
}

func staticFileContent(path string) ([]byte, error) {
	if initErr != nil {
		return nil, initErr
	}
	file, err := fileSys.Open(path)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil || fi.IsDir() {
		return nil, errors.New("not found")
	}
	defer file.Close()
	return fs.ReadFile(fileSys, path)
}
