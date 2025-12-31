package config

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-gorm/caches/v4"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/busyster996/dagflow/internal/pubsub"
	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/storage/models"
	"github.com/busyster996/dagflow/internal/utility"
	"github.com/busyster996/dagflow/internal/utility/sid"
	"github.com/busyster996/dagflow/pkg/logx"
)

var db *gorm.DB

func Init() error {
	viper.Set("log_dir", filepath.Join(viper.GetString("root_dir"), "logs"))
	viper.Set("script_dir", filepath.Join(viper.GetString("root_dir"), "scripts"))
	viper.Set("workspace_dir", filepath.Join(viper.GetString("root_dir"), "workspace"))
	if viper.GetString("node_name") == "" {
		viper.Set("node_name", fmt.Sprintf("%d-%d", viper.GetInt64("kind_id"), viper.GetInt64("node_id")))
	}

	var logfile string
	if viper.GetString("log_output") == "file" {
		logfile = filepath.Join(viper.GetString("log_dir"), utility.ServiceName+".log")
	}

	logx.SetupConsoleLogger(logfile, zap.AddStacktrace(zapcore.FatalLevel))
	level, err := zapcore.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		return fmt.Errorf("invalid log level: %v", err)
	}
	logx.SetLevel(level)

	var dirs = map[string]string{
		"root":      viper.GetString("root_dir"),
		"log":       viper.GetString("log_dir"),
		"script":    viper.GetString("script_dir"),
		"workspace": viper.GetString("workspace_dir"),
	}

	for name, dir := range dirs {
		if name == "log" && viper.GetString("log_output") != "file" {
			continue
		}
		if err = utility.EnsureDirExist(dir); err != nil {
			return fmt.Errorf("failed to ensure directory %s: %v", dir, err)
		}
		logx.Infof("%s dir: %s", name, dir)
	}

	logx.Infof("kind_id=%d node_id=%d", viper.GetInt64("kind_id"), viper.GetInt64("node_id"))

	// init pubsub queue
	if err = initPubsub(); err != nil {
		logx.Errorln(err)
		return err
	}

	if err = initStorage(); err != nil {
		logx.Errorln(err)
		return err
	}

	if viper.GetBool("enable_self_update") {
		utility.StartSelfUpdate(viper.GetString("self_url"), func() bool {
			if (storage.TaskCount(models.StateRunning) + storage.TaskCount(models.StatePending)) != 0 {
				// 还有任务执行中或者等待执行不升级
				logx.Warnln("the task has not been completed")
				return true
			}
			return false
		})
	}
	return nil
}

func initPubsub() error {
	if err := pubsub.New(viper.GetString("mq_url")); err != nil {
		logx.Errorln(err)
		return fmt.Errorf("failed to setup queue: %v", err)
	}
	return nil
}

func initStorage() (err error) {
	before, after, found := strings.Cut(viper.GetString("db_url"), "://")
	if !found {
		return errors.New("invalid storage url")
	}
	var dialector gorm.Dialector
	switch before {
	case storage.TypeMysql:
		// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
		dialector = mysql.Open(after)
	case storage.TypePostgres:
		// postgres://username:password@localhost:5432/database_name
		dialector = postgres.Open(viper.GetString("db_url"))
	case storage.TypeSqlserver:
		// sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm
		dialector = sqlserver.Open(viper.GetString("db_url"))
	case storage.TypeSqlite:
		dir := filepath.Join(viper.GetString("root_dir"), "data")
		if err = utility.EnsureDirExist(dir); err != nil {
			return fmt.Errorf("failed to ensure directory %s: %v", dir, err)
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.db3", utility.ServiceName))
		logx.Infof("%s file: %s", "data", file)
		viper.Set("db_url", fmt.Sprintf("%s://%s", storage.TypeSqlite, file))
		// test.db
		dialector = sqlite.Open(file)
	default:
		return errors.New("unsupported storage type")
	}
	gromConf := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable:       true,
			NoLowerCase:         false,
			IdentifierMaxLength: 256,
		},
		Logger: logger.New(logx.GetSubLogger(), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  logger.Error,
		}),
		SkipDefaultTransaction: true,
		FullSaveAssociations:   true,
		TranslateError:         true,
	}
	db, err = gorm.Open(dialector, gromConf)
	if err != nil {
		logx.Errorln(err)
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logx.Errorln(err)
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	_ = db.Use(&caches.Caches{Conf: &caches.Config{
		Easer: true,
	}})
	err = sid.Set(viper.GetUint64("kind_id"), viper.GetUint64("node_id"))
	if err != nil {
		logx.Errorln(err)
		return err
	}

	if err = storage.New(db); err != nil {
		logx.Errorln(err)
		return err
	}
	return
}

func GetGormDB() *gorm.DB {
	return db
}

func GetDB() (*sql.DB, error) {
	return db.DB()
}

func ClosePubsub(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()
	logx.Infoln("shutdown queue")
	pubsub.Shutdown(ctx)
}
func CloseDB() error {
	sqlDB, err := db.DB()
	if err != nil {
		logx.Errorln(err)
		return err
	}
	return sqlDB.Close()
}
