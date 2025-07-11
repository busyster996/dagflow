package worker

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/busyster996/dagflow/pkg/logx"
	"github.com/kardianos/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/busyster996/dagflow/config"
	"github.com/busyster996/dagflow/utils"
	"github.com/busyster996/dagflow/worker"
)

func New() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "worker",
		Short: "start a worker service",
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

			if strings.HasPrefix(strings.ToLower(viper.GetString("mq_url")), "inmemory") ||
				strings.HasPrefix(strings.ToLower(viper.GetString("db_url")), "sqlite") {
				return errors.New("memory queues or sqlite database is not allowed")
			}

			if viper.GetBool("enable_self_update") {
				utils.StartSelfUpdate(viper.GetString("self_url"))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := filepath.Abs(os.Args[0])
			if err != nil {
				logx.Errorln(err)
				return err
			}
			svc, err := service.New(&workerService{cmd.Context()}, &service.Config{
				Name:        utils.ServiceName,
				DisplayName: utils.ServiceName,
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
	cmd.PersistentFlags().String("node_name", "dagflow01", "node name")
	cmd.Flags().Int("pool_size", runtime.NumCPU()*2, "set the size of the execution work pool.")

	return cmd
}

type workerService struct {
	ctx context.Context
}

func (w *workerService) Start(s service.Service) error {
	// 调整工作池的大小
	worker.SetSize(viper.GetInt("pool_size"))
	return worker.Start(
		w.ctx,
		viper.GetString("node_name"),
		viper.GetString("work_space"),
		viper.GetString("script_dir"),
	)
}

func (w *workerService) Stop(s service.Service) error {
	// close pubsub
	config.ClosePubsub(w.ctx)
	// close worker
	worker.Shutdown()
	// close db
	if err := config.CloseDB(); err != nil {
		logx.Warnln(err)
	}
	logx.CloseLogger()
	return nil
}
