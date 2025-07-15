package workspace

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/internal/utils"
	"github.com/busyster996/dagflow/pkg/logx"
)

// Post
// @Summary		上传
// @Description	上传文件或目录
// @Tags		工作目录
// @Accept		multipart/form-data
// @Produce		application/json
// @Param		task path string true "任务名称"
// @Param		path query string false "路径"
// @Param		files formData file true "文件"
// @Success		200 {object} types.SBase[any]
// @Failure		500 {object} types.SBase[any]
// @Router		/api/v1/task/{task}/workspace [post]
func Post(c *gin.Context) {
	task := c.Param("task")
	if task == "" {
		base.Send(c, base.WithCode[any](types.CodeNoData).WithError(errors.New("task does not exist")))
		return
	}
	prefix := filepath.Join(viper.GetString("workspace_dir"), task)
	if !utils.FileOrPathExist(prefix) {
		base.Send(c, base.WithCode[any](types.CodeNoData).WithError(errors.New("task does not exist")))
		return
	}
	path := filepath.Join(prefix, utils.PathEscape(c.Query("path")))
	if err := base.SaveFiles(c, path); err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](types.CodeFailed).WithError(err))
		return
	}
	base.Send(c, base.WithCode[any](types.CodeSuccess))
}
