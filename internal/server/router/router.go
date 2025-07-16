package router

import (
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/busyster996/dagflow/internal/server/api/v1/event"
	"github.com/busyster996/dagflow/internal/server/api/v1/files"
	"github.com/busyster996/dagflow/internal/server/api/v1/pipeline"
	"github.com/busyster996/dagflow/internal/server/api/v1/pipeline/build"
	"github.com/busyster996/dagflow/internal/server/api/v1/task"
	"github.com/busyster996/dagflow/internal/server/api/v1/task/step"
	"github.com/busyster996/dagflow/internal/server/api/v1/task/workspace"
	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/router/middleware/zap"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/pkg/info"
	"github.com/busyster996/dagflow/pkg/logx"
)

// New
// @title			DAG Flow
// @version			1.0
// @Description		An `API` for cross-platform custom orchestration of execution steps
// @Description		without any third-party dependencies. Based on `DAG`, it implements the scheduling
// @Description		function of sequential execution of dependent steps and concurrent execution of
// @Description		non-dependent steps. <br /><br /> It provides `API` remote operation mode, batch
// @Description		execution of `Shell` , `Powershell` , `Python` and other commands, and easily
// @Description		completes common management tasks such as running automated operation and maintenance
// @Description		scripts, polling processes, installing or uninstalling software, updating applications,
// @Description		and installing patches.
// @contact.name	dagflow
// @contact.url		https://github.com/busyster996/dagflow/issues
// @license.name	GPL-3.0
// @license.url		https://github.com/busyster996/dagflow/blob/main/LICENSE
func New(gdb *gorm.DB, uploadDir string) (*gin.Engine, error) {
	router := gin.New()
	router.Use(
		zap.Logger,
		zap.Recovery,
		cors.Default(),
		gzip.Gzip(gzip.DefaultCompression),
		func(c *gin.Context) {
			c.Header("Server", "Gin")
			c.Header("X-Server", "Gin")
			c.Header("X-Version", info.Version)
			c.Header("X-Powered-By", info.UserEmail)
		},
	)
	// debug pprof
	pprof.Register(router)
	// base
	router.GET("/version", version)
	router.GET("/healthyz", healthyz)
	router.GET("/heartbeat", heartbeat)
	router.HEAD("/heartbeat", heartbeat)
	// swagger
	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiV1 := router.Group("/api/v1")
	// V1
	{
		// event
		apiV1.GET("/event", event.Stream)

		// pipeline
		apiV1.GET("/pipeline", pipeline.List)
		apiV1.POST("/pipeline", pipeline.Post)
		apiV1.POST("/pipeline/:pipeline", pipeline.Update)
		apiV1.GET("/pipeline/:pipeline", pipeline.Detail)
		apiV1.DELETE("/pipeline/:pipeline", pipeline.Delete)

		// build
		apiV1.GET("/pipeline/:pipeline/build", build.List)
		apiV1.POST("/pipeline/:pipeline/build", build.Post)
		apiV1.GET("/pipeline/:pipeline/build/:build", build.Detail)
		apiV1.POST("/pipeline/:pipeline/build/:build", build.ReRun)
		apiV1.DELETE("/pipeline/:pipeline/build/:build", build.Delete)

		// task
		apiV1.GET("/task", task.List)
		apiV1.POST("/task", task.Post)
		apiV1.GET("/task/:task", task.Detail)
		apiV1.PUT("/task/:task", task.Manager)
		apiV1.DELETE("/task/:task", task.Delete)
		apiV1.GET("/task/:task/dump", task.Dump)

		// workspace
		apiV1.GET("/task/:task/workspace", workspace.Get)
		apiV1.DELETE("/task/:task/workspace", workspace.Delete)
		apiV1.POST("/task/:task/workspace", workspace.Post)

		// step
		apiV1.GET("/task/:task/step", step.List)
		apiV1.GET("/task/:task/step/:step", step.Detail)
		apiV1.PUT("/task/:task/step/:step", step.Manager)
		apiV1.GET("/task/:task/step/:step/log", step.Log)

		// tus file server
		if err := files.New(uploadDir, path.Join(apiV1.BasePath(), "/files/"), gdb); err != nil {
			logx.Errorln(err)
			return nil, err
		}

		apiV1.Any("/files", files.Handler())
		apiV1.Any("/files/*any", files.Handler())
	}

	// no method
	router.NoMethod(func(c *gin.Context) {
		base.Send(c, base.WithCode[any](types.CodeNoData).WithError(errors.New("method not allowed")))
	})

	// no route
	router.NoRoute(staticHandler())
	return router, nil
}
