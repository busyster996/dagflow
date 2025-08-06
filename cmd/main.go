//go:build (windows || linux || darwin) && (amd64 || arm64 || 386)

package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/busyster996/dagflow/cmd/api"
	"github.com/busyster996/dagflow/cmd/standalone"
	"github.com/busyster996/dagflow/cmd/worker"
	"github.com/busyster996/dagflow/internal/utils"
	"github.com/busyster996/dagflow/pkg/info"
	"github.com/busyster996/dagflow/pkg/logx"
)

const longText = `An API for cross-platform custom orchestration of execution steps without any third-party dependencies. 
Based on DAG, it implements the scheduling function of sequential execution of dependent steps and concurrent execution of non-dependent steps.`

func init() {
	time.Local = time.UTC
	viper.SetOptions(
		viper.KeyDelimiter("::"),
		viper.ExperimentalBindStruct(),
	)
}

func main() {
	var cmd = &cobra.Command{
		Use:   os.Args[0],
		Short: "Operating system remote execution interface",
		Long:  longText,
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       info.Version,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.SetHelpTemplate(`{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)
	cmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	_ = cmd.PersistentFlags().MarkShorthandDeprecated("help", "please use --help")
	cmd.PersistentFlags().BoolP("version", "v", false, "Print version information and quit")

	cmd.PersistentFlags().Int64("node_id", 0, "node id")
	cmd.PersistentFlags().Int64("kind_id", 0, "kind id")
	cmd.PersistentFlags().String("root_dir", utils.DefaultDir(), "root directory")
	cmd.PersistentFlags().String("log_output", "file", "log output [file,stdout]")
	cmd.PersistentFlags().String("log_level", "debug", "log level [debug,info,warn,error]")
	cmd.PersistentFlags().Bool("enable_self_update", true, "enable self update")
	cmd.PersistentFlags().String("self_url", "https://oss.yfdou.com/tools/dagflow", "self Update URL")
	cmd.PersistentFlags().String("mq_url", "inmemory://localhost", "message queue url. [inmemory,amqp]")
	cmd.PersistentFlags().String("db_url", "sqlite://localhost", "database type. [sqlite,mysql,postgres,sqlserver]")

	cmd.AddCommand(
		standalone.New(),
		api.New(),
		worker.New(),
		&cobra.Command{
			Use:   "version",
			Short: "print version information and quit",
			RunE: func(cmd *cobra.Command, args []string) error {
				info.PrintHeadInfo()
				return nil
			},
		},
	)

	if err := cmd.Execute(); err != nil {
		logx.Fatalln(err)
	}
}
