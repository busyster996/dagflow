package api

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kardianos/service"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/busyster996/dagflow/internal/config"
	"github.com/busyster996/dagflow/internal/server"
	"github.com/busyster996/dagflow/internal/utility"
	"github.com/busyster996/dagflow/pkg/logx"
)

func New() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "api",
		Short: "start api service",
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

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name, err := filepath.Abs(os.Args[0])
			if err != nil {
				logx.Errorln(err)
				return err
			}
			svc, err := service.New(&apiService{cmd.Context()}, &service.Config{
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
	cmd.Flags().String("addr", "tcp://0.0.0.0:2376", "listening address.")
	cmd.Flags().Duration("exec_timeout", 24*time.Hour, "set the task exec command expire time")

	return cmd
}

type apiService struct {
	ctx context.Context
}

func (a *apiService) Start(s service.Service) error {
	return server.Start(a.ctx)
}

func (a *apiService) Stop(s service.Service) error {
	// close http server
	if err := server.Close(); err != nil {
		logx.Warnln(err)
	}
	// close pubsub
	config.ClosePubsub(a.ctx)
	// close db
	if err := config.CloseDB(); err != nil {
		logx.Warnln(err)
	}
	logx.CloseLogger()
	return nil
}
