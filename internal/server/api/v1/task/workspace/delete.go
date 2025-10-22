package workspace

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/utility"
	"github.com/busyster996/dagflow/pkg/logx"
)

// Delete
// @Summary		删除
// @Description	删除指定目录或文件
// @Tags		工作目录
// @Accept		application/json
// @Produce		application/json
// @Param		task path string true "任务名称"
// @Param		path query string false "路径"
// @Success		200 {object} base.IResponse[any]
// @Failure		500 {object} base.IResponse[any]
// @Router		/api/v1/task/{task}/workspace [delete]
func Delete(c *gin.Context) {
	task := c.Param("task")
	if task == "" {
		base.Send(c, base.WithCode[any](base.CodeNoData).WithError(errors.New("task does not exist")))
		return
	}
	prefix := filepath.Join(viper.GetString("workspace_dir"), task)
	if !utility.FileOrPathExist(prefix) {
		base.Send(c, base.WithCode[any](base.CodeNoData).WithError(errors.New("task does not exist")))
		return
	}
	path := filepath.Join(prefix, utility.PathEscape(c.Query("path")))
	if err := os.RemoveAll(path); err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](base.CodeFailed).WithError(err))
		return
	}
	base.Send(c, base.WithCode[any](base.CodeSuccess))
}
