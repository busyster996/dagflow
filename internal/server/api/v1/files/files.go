package files

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"

	"github.com/busyster996/dagflow/internal/utils"
	"github.com/busyster996/dagflow/pkg/logx"
	"github.com/busyster996/dagflow/pkg/tus"
	filelocker "github.com/busyster996/dagflow/pkg/tus/locker/file"
	filestore "github.com/busyster996/dagflow/pkg/tus/storage/file"
	"github.com/busyster996/dagflow/pkg/tus/types"
)

var tunServer *tus.STusx

func New(uploadDir, basePath string, gdb *gorm.DB) error {
	if tunServer != nil {
		logx.Errorln("tus server already started")
		return nil
	}
	var locker = filelocker.New(filepath.Join(uploadDir, ".lock"))
	store, err := filestore.New(filepath.Join(os.TempDir(), ".tusd"), gdb, locker)
	if err != nil {
		logx.Errorln(err)
		return err
	}
	store.Cleanup(context.Background(), 1*time.Hour)
	tunServer, err = tus.New(&tus.SConfig{
		BasePath: basePath,
		Store:    store,
		Logger:   logx.GetSubLogger(),
		PreUploadCreateCallback: func(hook types.HookEvent) (types.HTTPResponse, types.FileInfoChanges, error) {
			id := ksuid.New().String()
			taskID, ok := hook.Upload.MetaData["task_id"]
			if !ok {
				return types.HTTPResponse{
					StatusCode: http.StatusBadRequest,
					Body:       "task_id is required",
				}, types.FileInfoChanges{}, errors.New("task_id is required")
			}

			//if !service.Task(taskID).Exist() {
			//	return types.HTTPResponse{
			//		StatusCode: http.StatusBadRequest,
			//		Body:       "taskid is not exist",
			//	}, types.FileInfoChanges{}, errors.New("taskid is not exist")
			//}

			return types.HTTPResponse{}, types.FileInfoChanges{
				ID: filepath.Join(taskID, id),
			}, nil
		},
		PreFinishResponseCallback: func(hook types.HookEvent) (types.HTTPResponse, error) {
			if hook.Upload.IsFinal {
				filename := hook.Upload.MetaData["filename"]
				if filename == "" {
					filename = filepath.Base(hook.Upload.ID)
				}

				src := filepath.Join(os.TempDir(), ".tusd", hook.Upload.ID)
				dst := filepath.Join(uploadDir, filepath.Dir(hook.Upload.ID), filename)
				if err = utils.CopyFile(src, dst); err != nil {
					return types.HTTPResponse{
						StatusCode: http.StatusInternalServerError,
						Body:       "failed to copy file",
					}, err
				}
			}
			return types.HTTPResponse{
				Headers: map[string]string{
					"ID":   filepath.Base(hook.Upload.ID),
					"Path": filepath.Dir(hook.Upload.ID),
				},
			}, nil
		},
	})
	if err != nil {
		logx.Errorln(err)
		return err
	}
	return nil
}

func Handler() gin.HandlerFunc {
	return gin.WrapH(tunServer)
}
