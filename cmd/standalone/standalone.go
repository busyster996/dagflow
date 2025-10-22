package standalone

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/busyster996/dagflow/internal/config"
	"github.com/busyster996/dagflow/internal/server"
	"github.com/busyster996/dagflow/internal/utility"
	"github.com/busyster996/dagflow/internal/worker"
	"github.com/busyster996/dagflow/pkg/logx"
)

func New() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "standalone",
		Short: "start standalone service",
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			_ = viper.BindPFlags(cmd.PersistentFlags())
			_ = viper.BindPFlags(cmd.Flags())

			if err := config.Init(); err != nil {
				logx.Errorln(err)
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := filepath.Abs(os.Args[0])
			if err != nil {
				logx.Errorln(err)
				return err
			}
			svc, err := service.New(&standaloneService{cmd.Context()}, &service.Config{
				Name:        utility.ServiceName,
				DisplayName: utility.ServiceName,
				Description: "Operating System Remote Executor Api",
				Executable:  name,
				Arguments:   os.Args[1:],
			})
			if err != nil {
				return err
			}
			err = svc.Run()
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().String("relative_path", "/", "web relative path")
	cmd.Flags().String("addr", "tcp://0.0.0.0:2376", "listening address.")
	cmd.Flags().Duration("exec_timeout", 24*time.Hour, "set the task exec command expire time")

	cmd.Flags().String("node_name", "dagflow01", "node name")
	cmd.Flags().Int("pool_size", runtime.NumCPU()*2, "set the size of the execution work pool.")
	return cmd
}

type standaloneService struct {
	ctx context.Context
}

func (a *standaloneService) Start(s service.Service) error {
	// 调整工作池的大小
	worker.SetSize(viper.GetInt("pool_size"))
	if err := worker.Start(a.ctx); err != nil {
		return err
	}
	return server.Start(a.ctx)
}

func (a *standaloneService) Stop(s service.Service) error {
	// close http server
	if err := server.Close(); err != nil {
		logx.Warnln(err)
	}
	// close pubsub
	config.ClosePubsub(a.ctx)
	// close worker
	worker.Shutdown()
	// close db
	if err := config.CloseDB(); err != nil {
		logx.Warnln(err)
	}
	logx.CloseLogger()
	return nil
}
