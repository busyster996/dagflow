package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/busyster996/dagflow/pkg/tus/locker"
	"github.com/busyster996/dagflow/pkg/tus/storage"
	"github.com/busyster996/dagflow/pkg/tus/types"
)

var (
	defaultFilePerm      = os.FileMode(0664)
	defaultDirectoryPerm = os.FileMode(0754)
	bufferPool           = sync.Pool{
		New: func() interface{} {
			return make([]byte, 64*1024*1024) // 64MB缓冲区
		},
	}
)

// FileUploadChunks GORM模型定义
type FileUploadChunks struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	FileID       string         `gorm:"primaryKey;uniqueIndex;size:255;comment:文件ID" json:"file_id"`
	FileSize     int64          `gorm:"not null;comment:文件大小" json:"file_size"`
	OffsetSize   int64          `gorm:"not null;default:0;comment:偏移量" json:"offset_size"`
	IsPartial    bool           `gorm:"default:false;comment:是否为分片" json:"is_partial"`
	MetadataInfo datatypes.JSON `gorm:"type:json;comment:元数据" json:"metadata_info"`
	PartialIDs   datatypes.JSON `gorm:"type:json;comment:分片ID" json:"partial_ids"`
}

// TableName 指定表名
func (FileUploadChunks) TableName() string {
	return "file_upload_chunks"
}

type SFileStore struct {
	Dir    string
	db     *gorm.DB
	locker locker.ILocker
}

func New(dir string, db *gorm.DB, locker locker.ILocker) (*SFileStore, error) {
	_ = os.MkdirAll(dir, defaultDirectoryPerm)

	store := &SFileStore{
		Dir:    dir,
		db:     db,
		locker: locker,
	}

	// 配置GORM
	if err := store.configureGORM(); err != nil {
		return nil, fmt.Errorf("failed to configure GORM: %w", err)
	}

	// 自动迁移
	if err := store.autoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return store, nil
}

// 配置GORM
func (store *SFileStore) configureGORM() error {
	if store.db.Dialector.Name() == "sqlite" {
		// SQLite特殊配置
		optimizations := []string{
			"PRAGMA mode=rwc;",
			"PRAGMA busy_timeout = 60000;",
			"PRAGMA journal_mode = WAL;",
			"PRAGMA synchronous = NORMAL;",
			"PRAGMA cache = shared;",
			"PRAGMA cache_spill = ON;",
			"PRAGMA cache_size = -131072;",
			"PRAGMA foreign_keys = ON;",
			"PRAGMA temp_store = MEMORY;",
			"PRAGMA mmap_size = 536870912;",
			"PRAGMA wal_autocheckpoint = 1000;",
			"PRAGMA locking_mode = NORMAL;",
			"PRAGMA read_uncommitted = ON;",
			"PRAGMA journal_size_limit=104857600;",
		}

		for _, sqlStr := range optimizations {
			if err := store.db.Exec(sqlStr).Error; err != nil {
				fmt.Printf("Warning: failed to execute %s: %v\n", sqlStr, err)
			}
		}

	}

	return nil
}

func (store *SFileStore) autoMigrate() error {
	return store.db.AutoMigrate(&FileUploadChunks{})
}

func (store *SFileStore) binPath(id string) string {
	return filepath.Join(store.Dir, id)
}

func (store *SFileStore) NewUpload(ctx context.Context, info types.FileInfo) (storage.IUpload, error) {
	if info.ID == "" {
		info.ID = types.Uid()
	}

	upload := &sFileUpload{
		info:    info,
		binPath: store.binPath(info.ID),
		store:   store,
	}

	binLock, err := store.locker.NewLock(strings.ReplaceAll(strings.TrimSpace(upload.binPath), "/", ":"))
	if err != nil {
		return nil, err
	}
	upload.binLock = binLock

	if err = upload.binLock.Lock(ctx); err != nil {
		return nil, err
	}
	defer upload.binLock.Unlock()

	if err = upload.createFile(upload.binPath, nil); err != nil {
		return nil, err
	}

	if err = upload.writeInfo(ctx); err != nil {
		return nil, err
	}

	return upload, nil
}

func (store *SFileStore) GetUpload(ctx context.Context, id string) (storage.IUpload, error) {
	upload := &sFileUpload{
		binPath: store.binPath(id),
		store:   store,
	}

	binLock, err := store.locker.NewLock(strings.ReplaceAll(strings.TrimSpace(upload.binPath), "/", ":"))
	if err != nil {
		return nil, err
	}
	upload.binLock = binLock

	if err = upload.readInfo(ctx, id); err != nil {
		return nil, err
	}

	stat, err := os.Stat(upload.binPath)
	if err != nil {
		return nil, err
	}
	upload.info.Offset = stat.Size()
	if upload.info.Offset != upload.info.Offset {
		if err = upload.updateOffset(ctx); err != nil {
			return nil, err
		}
	}

	return upload, nil
}

func (store *SFileStore) Cleanup(ctx context.Context, expiredBefore time.Duration) {
	go func() {
		// 定时清理
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				store.cleanup(ctx, expiredBefore)
			}
		}
	}()
}

func (store *SFileStore) cleanup(ctx context.Context, expiredBefore time.Duration) {
	lock, err := store.locker.NewLock("cleanup")
	if err != nil {
		fmt.Printf("failed to get cleanup lock: %v\n", err)
		return
	}
	if err = lock.Lock(ctx); err != nil {
		fmt.Printf("failed to get cleanup lock: %v\n", err)
		return
	}
	defer lock.Unlock()
	var (
		expiredTime = time.Now().Add(-expiredBefore)
		uploadIDs   []string
	)

	result := store.db.WithContext(ctx).
		Model(&FileUploadChunks{}).
		Select("file_id").
		Where("created_at < ?", expiredTime).
		Find(&uploadIDs)
	if result.Error != nil {
		fmt.Printf("failed to get expired uploads: %v\n", result.Error)
		return
	}

	for _, uploadID := range uploadIDs {
		err = os.RemoveAll(store.binPath(uploadID))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			fmt.Printf("failed to remove expired upload: %v\n", err)
			continue
		}
		store.db.WithContext(ctx).Where("file_id = ?", uploadID).Delete(&FileUploadChunks{})
	}
}

type sFileUpload struct {
	binLock locker.ILock
	info    types.FileInfo
	binPath string
	store   *SFileStore
}

func (upload *sFileUpload) writeInfo(ctx context.Context) error {
	var (
		metadata   []byte
		partialIDs []byte
	)

	if len(upload.info.PartialIDs) > 0 {
		var err error
		metadata, err = json.Marshal(upload.info.MetaData)
		if err != nil {
			return err
		}
	}
	if len(upload.info.PartialIDs) > 0 {
		var err error
		partialIDs, err = json.Marshal(upload.info.PartialIDs)
		if err != nil {
			return err
		}
	}
	info := &FileUploadChunks{
		FileID:       upload.info.ID,
		FileSize:     upload.info.Size,
		OffsetSize:   upload.info.Offset,
		IsPartial:    upload.info.IsPartial,
		MetadataInfo: datatypes.JSON(metadata),
		PartialIDs:   datatypes.JSON(partialIDs),
	}
	var doUpdates = []string{
		"file_size",
		"offset_size",
		"is_partial",
	}
	if metadata != nil {
		doUpdates = append(doUpdates, "metadata_info")
	}
	if partialIDs != nil {
		doUpdates = append(doUpdates, "partial_ids")
	}

	result := upload.store.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "file_id"}},
			DoUpdates: clause.AssignmentColumns(doUpdates),
		}).Create(info)
	return result.Error
}

func (upload *sFileUpload) readInfo(ctx context.Context, id string) error {
	var info FileUploadChunks
	result := upload.store.db.WithContext(ctx).Where("file_id = ?", id).First(&info)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("upload not found")
		}
		return result.Error
	}
	upload.info.ID = info.FileID
	upload.info.Size = info.FileSize
	upload.info.Offset = info.OffsetSize
	upload.info.IsPartial = info.IsPartial

	if len(info.MetadataInfo) > 0 {
		if err := json.Unmarshal(info.MetadataInfo, &upload.info.MetaData); err != nil {
			return err
		}
	}
	if len(info.PartialIDs) > 0 {
		if err := json.Unmarshal(info.PartialIDs, &upload.info.PartialIDs); err != nil {
			return err
		}
	}
	return nil
}

func (upload *sFileUpload) updateOffset(ctx context.Context) error {
	return upload.store.db.WithContext(ctx).Model(&FileUploadChunks{}).
		Where("file_id = ?", upload.info.ID).
		Update("offset_size", upload.info.Offset).Error
}

func (upload *sFileUpload) createFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), defaultDirectoryPerm); err != nil {
		return fmt.Errorf("failed to create directory for %s: %s", path, err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, defaultFilePerm)
	if err != nil {
		return err
	}
	if content != nil {
		if _, err = file.Write(content); err != nil {
			return err
		}
	}
	return file.Close()
}

func (upload *sFileUpload) GetInfo(ctx context.Context) (types.FileInfo, error) {
	if err := upload.readInfo(ctx, upload.info.ID); err != nil {
		return types.FileInfo{}, err
	}
	stat, err := os.Stat(upload.binPath)
	if err != nil {
		return types.FileInfo{}, fmt.Errorf("upload not found")
	}
	upload.info.Offset = stat.Size()
	return upload.info, nil
}

func (upload *sFileUpload) GetReader(ctx context.Context) (io.ReadCloser, error) {
	return os.Open(upload.binPath)
}

func (upload *sFileUpload) WriteChunk(ctx context.Context, offset int64, src io.Reader) (int64, error) {
	if err := upload.binLock.Lock(ctx); err != nil {
		return 0, err
	}
	defer upload.binLock.Unlock()

	file, err := os.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close()
	}()

	if _, err = file.Seek(offset, io.SeekStart); err != nil {
		return 0, err
	}

	buffer := bufferPool.Get().([]byte)
	defer bufferPool.Put(buffer)

	n, err := io.CopyBuffer(file, src, buffer)
	if err != nil {
		return n, err
	}

	upload.info.Offset += n
	return n, upload.writeInfo(ctx)
}

func (upload *sFileUpload) ConcatUploads(ctx context.Context, uploads []storage.IUpload) (err error) {
	if err = upload.binLock.Lock(ctx); err != nil {
		return err
	}
	defer upload.binLock.Unlock()

	file, err := os.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	for _, partialUpload := range uploads {
		_partialUpload := partialUpload.(*sFileUpload)
		if err = _partialUpload.appendTo(ctx, file); err != nil {
			return err
		}
		if err = _partialUpload.Terminate(ctx); err != nil {
			return err
		}
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}
	upload.info.Size = info.Size()
	upload.info.Offset = info.Size()
	if err = upload.writeInfo(ctx); err != nil {
		return err
	}
	return
}

func (upload *sFileUpload) appendTo(ctx context.Context, file *os.File) error {
	if err := upload.binLock.Lock(ctx); err != nil {
		return err
	}
	defer upload.binLock.Unlock()

	src, err := os.Open(upload.binPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = src.Close()
	}()
	buffer := bufferPool.Get().([]byte)
	defer bufferPool.Put(buffer)

	_, err = io.CopyBuffer(file, src, buffer)
	return err
}

func (upload *sFileUpload) ServeContent(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if err := upload.binLock.Lock(ctx); err != nil {
		return err
	}
	defer upload.binLock.Unlock()
	http.ServeFile(w, r, upload.binPath)
	return nil
}

func (upload *sFileUpload) Terminate(ctx context.Context) error {
	if err := upload.binLock.Lock(ctx); err != nil {
		return err
	}
	defer upload.binLock.Unlock()

	err := upload.store.db.WithContext(ctx).Where("file_id = ?", upload.info.ID).Delete(&FileUploadChunks{}).Error
	if err != nil {
		return err
	}

	err = os.RemoveAll(upload.binPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}
