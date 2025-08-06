package router

import (
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/pkg/info"
)

// version
// @Summary		版本
// @Description	获取版本信息
// @Tags		默认
// @Accept		application/json
// @Produce		application/json
// @Success		200 {object} types.SBase[types.SVersion]
// @Failure		500 {object} types.SBase[any]
// @Router		/version [get]
func version(c *gin.Context) {
	base.Send(c, base.WithData(&types.SVersion{
		Version:   info.Version,
		BuildTime: info.BuildTime,
		Git: &types.SVersionGit{
			URL:    info.GitUrl,
			Branch: info.GitBranch,
			Commit: info.GitCommit,
		},
		Go: &types.SVersionGO{
			Version: runtime.Version(),
			OS:      runtime.GOOS,
			Arch:    runtime.GOARCH,
		},
		User: &types.SVersionUser{
			Name:  info.UserName,
			Email: info.UserEmail,
		},
	}))
}
