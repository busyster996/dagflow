package storage

import (
	"gorm.io/gorm"

	"github.com/busyster996/dagflow/internal/storage/models"
	"github.com/busyster996/dagflow/pkg/logx"
)

var storage IStorage

const (
	TypeSqlite    = "sqlite"
	TypeMysql     = "mysql"
	TypePostgres  = "postgres"
	TypeSqlserver = "sqlserver"
)

func New(gdb *gorm.DB) error {
	db := &sDatabase{DB: gdb}

	if gdb.Name() == TypeSqlite {
		logx.Infoln("init sqlite")
		db.initSqlite()
	}

	// 自动迁移表
	if err := db.AutoMigrate(
		&models.STask{},
		&models.STaskEnv{},
		&models.SStep{},
		&models.SStepEnv{},
		&models.SStepDepend{},
		&models.SStepLog{},
		&models.SPipeline{},
		&models.SPipelineBuild{},
	); err != nil {
		logx.Errorln(err)
		return err
	}
	storage = db
	return nil
}

func Name() string {
	return storage.Name()
}

func FixDatabase(node string) (err error) {
	return storage.FixDatabase(node)
}

func Task(name string) ITask {
	return storage.Task(name)
}

func NodeTasks(node string) []ITask {
	return storage.NodeTasks(node)
}

func TaskCreate(task *models.STask) (err error) {
	return storage.TaskCreate(task)
}

func TaskCount(state models.State) (res int64) {
	return storage.TaskCount(state)
}

func TaskList(page, pageSize int64, str string) (res []*models.STask, total int64) {
	return storage.TaskList(page, pageSize, str)
}

func Pipeline(name string) IPipeline {
	return storage.Pipeline(name)
}

func PipelineCreate(pipeline *models.SPipeline) (err error) {
	return storage.PipelineCreate(pipeline)
}

func PipelineList(page, pageSize int64, str string) (res []*models.SPipeline, total int64) {
	return storage.PipelineList(page, pageSize, str)
}
