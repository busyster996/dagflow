package storagev2

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-gorm/caches/v4"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/storagev2/model"
	"github.com/busyster996/dagflow/internal/utils"
	"github.com/busyster996/dagflow/internal/utils/sid"
	"github.com/busyster996/dagflow/pkg/logx"
)

var db *gorm.DB

const (
	TypeSqlite    = "sqlite"
	TypeMysql     = "mysql"
	TypePostgres  = "postgres"
	TypeSqlserver = "sqlserver"
)

func New() error {
	err := sid.Set(viper.GetInt64("kind_id"), viper.GetInt64("node_id"))
	if err != nil {
		logx.Errorln(err)
		return err
	}

	before, after, found := strings.Cut(viper.GetString("db_url"), "://")
	if !found {
		return errors.New("invalid storage url")
	}
	var dialector gorm.Dialector
	switch before {
	case TypeMysql:
		// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
		dialector = mysql.Open(after)
	case TypePostgres:
		// postgres://username:password@localhost:5432/database_name
		dialector = postgres.Open(viper.GetString("db_url"))
	case TypeSqlserver:
		// sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm
		dialector = sqlserver.Open(viper.GetString("db_url"))
	case TypeSqlite:
		dir := filepath.Join(viper.GetString("root_dir"), "data")
		if err := utils.EnsureDirExist(dir); err != nil {
			return fmt.Errorf("failed to ensure directory %s: %v", dir, err)
		}
		file := filepath.Join(dir, fmt.Sprintf("%s.db3", utils.ServiceName))
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
	err = db.Use(&caches.Caches{Conf: &caches.Config{
		Easer: true,
	}})
	if err != nil {
		logx.Errorln(err)
		return err
	}

	db, err = gorm.Open(dialector, gromConf)
	if err != nil {
		logx.Errorln(err)
		return err
	}

	// 自动迁移表
	if err = db.AutoMigrate(
		&model.Param{},
		&model.Task{},
		&model.TaskParam{},
		&model.Step{},
		&model.StepParam{},
		&model.StepDepend{},
		&model.Execution{},
		&model.Output{},
	); err != nil {
		logx.Errorln(err)
		return err
	}
	return initSqlite(before)
}

func initSqlite(driver string) error {
	if driver != TypeSqlite {
		return nil
	}
	var sqls = []string{
		"PRAGMA mode=rwc;",
		"PRAGMA foreign_keys=ON;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA journal_mode=WAL;",
		"PRAGMA journal_size_limit=104857600;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA cache=shared;",
		"PRAGMA cache_size=-8000;",
		"PRAGMA mmap_size=134217728;",
		"PRAGMA temp_store=MEMORY;",
		"PRAGMA locking_mode=NORMAL;",
		"PRAGMA cache_spill=ON;",
	}
	return db.Transaction(func(tx *gorm.DB) error {
		for _, sql := range sqls {
			if err := tx.Exec(sql).Error; err != nil {
				return err // Return error if any PRAGMA fails
			}
		}
		return nil
	})
}
