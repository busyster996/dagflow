package build

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/service"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/pkg/logx"
)

// Post
// @Summary 	创建
// @Description 创建构建任务
// @Tags 		构建
// @Accept		application/json
// @Produce		application/json
// @Param		pipeline path string true "流水线名称"
// @Param		build body types.SPipelineBuildReq true "构建参数"
// @Success		200 {object} base.IResponse[types.STaskCreateRes]
// @Failure		500 {object} base.IResponse[any]
// @Router		/api/v1/pipeline/{pipeline}/build [post]
func Post(c *gin.Context) {
	pipelineName := c.Param("pipeline")
	if pipelineName == "" {
		base.Send(c, base.WithCode[any](base.CodeNoData).WithError(errors.New("pipeline does not exist")))
		return
	}
	var req = &types.SPipelineBuildReq{
		Params: make(map[string]any),
	}
	if err := c.ShouldBind(req); err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](base.CodeFailed).WithError(err))
		return
	}
	build, err := service.Pipeline(pipelineName).BuildCreate(req)
	if err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](base.CodeFailed).WithError(err))
		return
	}

	c.Request.Header.Set(types.XTaskName, build)
	c.Header(types.XTaskName, build)

	base.Send(c, base.WithData(&types.STaskCreateRes{
		Name: build,
	}))
}

// ReRun
// @Summary 	重新运行
// @Description 重新执行指定构建任务
// @Tags 		构建
// @Accept		application/json
// @Produce		application/json
// @Param		pipeline path string true "流水线名称"
// @Param		build path string true "构建名称"
// @Success		200 {object} base.IResponse[any]
// @Failure		500 {object} base.IResponse[any]
// @Router		/api/v1/pipeline/{pipeline}/build/{build} [post]
func ReRun(c *gin.Context) {
	pipelineName := c.Param("pipeline")
	if pipelineName == "" {
		base.Send(c, base.WithCode[any](base.CodeNoData).WithError(errors.New("pipeline does not exist")))
		return
	}
	buildName := c.Param("build")
	if buildName == "" {
		base.Send(c, base.WithCode[any](base.CodeNoData).WithError(errors.New("build does not exist")))
		return
	}
	err := service.Pipeline(pipelineName).BuildReRun(buildName)
	if err != nil {
		logx.Errorln(err)
		base.Send(c, base.WithCode[any](base.CodeFailed).WithError(err))
		return
	}
	base.Send(c, base.WithCode[any](base.CodeSuccess))
}
