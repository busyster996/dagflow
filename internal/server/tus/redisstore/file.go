// Package redisstore provide a storage backend based on Redis for metadata and local file system for binary data.
//
// FileStore is a storage backend used as a handler.DataStore in handler.NewHandler.
// It stores the upload metadata in Redis and binary data in local files.
// The binary files are stored as `[id]` files without an extension containing the raw uploaded data.
// Upload metadata is stored in Redis with automatic TTL cleanup.
package redisstore

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tus/tusd/v2/pkg/handler"
)

var defaultFilePerm = os.FileMode(0664)
var defaultDirectoryPerm = os.FileMode(0754)

const (
	// StorageKeyPath is the key of the path of uploaded file in handler.FileInfo.Storage
	StorageKeyPath = "Path"
	// StorageKeyRedisKey is the key of the Redis key in handler.FileInfo.Storage
	StorageKeyRedisKey = "RedisKey"

	// Redis相关常量
	InfoKeyPrefix = "tusd:upload:info:" // Redis key前缀
	DefaultTTL    = 2 * time.Hour       // 默认过期时间2小时
)

// FileStore combines Redis for metadata storage and local filesystem for binary data
type FileStore struct {
	// Relative or absolute path to store binary files in
	Path string
	// Redis client for storing upload metadata
	RedisClient *redis.Client
}

// NewFromClient creates a new file based storage backend with Redis metadata storage
func NewFromClient(path string, redisClient *redis.Client) *FileStore {
	return &FileStore{
		Path:        path,
		RedisClient: redisClient,
	}
}

// New creates a new file based storage backend
func New(path string, uri string) (*FileStore, error) {
	connection, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(connection)
	if res := client.Ping(context.Background()); res.Err() != nil {
		return nil, res.Err()
	}
	return &FileStore{
		Path:        path,
		RedisClient: client,
	}, nil
}

// UseIn sets this store as the core data store in the passed composer and adds
// all possible extension to it.
func (store FileStore) UseIn(composer *handler.StoreComposer) {
	composer.UseCore(store)
	composer.UseTerminater(store)
	composer.UseConcater(store)
	composer.UseLengthDeferrer(store)
	composer.UseContentServer(store)
}

func (store FileStore) NewUpload(ctx context.Context, info handler.FileInfo) (handler.Upload, error) {
	if info.ID == "" {
		info.ID = uid()
	}

	// Redis key for storing upload metadata
	redisKey := InfoKeyPrefix + info.ID

	// The binary file's location might be modified by the pre-create hook
	var binPath string
	if info.Storage != nil && info.Storage[StorageKeyPath] != "" {
		// filepath.Join treats absolute and relative paths the same, so we must
		// handle them on our own. Absolute paths get used as-is, while relative
		// paths are joined to the storage path.
		if filepath.IsAbs(info.Storage[StorageKeyPath]) {
			binPath = info.Storage[StorageKeyPath]
		} else {
			binPath = filepath.Join(store.Path, info.Storage[StorageKeyPath])
		}
	} else {
		binPath = store.defaultBinPath(info.ID)
	}

	info.Storage = map[string]string{
		"Type":             "filestore-redis",
		StorageKeyPath:     binPath,
		StorageKeyRedisKey: redisKey,
	}

	// Create binary file with no content
	if err := createFile(binPath, nil); err != nil {
		return nil, fmt.Errorf("failed to create binary file: %w", err)
	}

	upload := &fileUpload{
		info:        info,
		redisKey:    redisKey,
		binPath:     binPath,
		redisClient: store.RedisClient,
	}

	// Store metadata in Redis
	if err := upload.writeInfo(ctx); err != nil {
		// Clean up binary file if Redis write fails
		_ = os.Remove(binPath)
		return nil, fmt.Errorf("failed to store upload metadata: %w", err)
	}

	return upload, nil
}

func (store FileStore) GetUpload(ctx context.Context, id string) (handler.Upload, error) {
	redisKey := InfoKeyPrefix + id

	// Get upload metadata from Redis
	data, err := store.RedisClient.Get(ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, handler.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get upload metadata from Redis: %w", err)
	}

	var info handler.FileInfo
	if err := json.Unmarshal([]byte(data), &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal upload metadata: %w", err)
	}

	// Get binary file path from metadata
	var binPath string
	if info.Storage != nil && info.Storage[StorageKeyPath] != "" {
		binPath = info.Storage[StorageKeyPath]
	} else {
		binPath = store.defaultBinPath(info.ID)
	}

	// Check if binary file exists and get current size
	stat, err := os.Stat(binPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, handler.ErrNotFound
		}
		return nil, fmt.Errorf("failed to stat binary file: %w", err)
	}

	info.Offset = stat.Size()

	return &fileUpload{
		info:        info,
		binPath:     binPath,
		redisKey:    redisKey,
		redisClient: store.RedisClient,
	}, nil
}

func (store FileStore) AsTerminatableUpload(upload handler.Upload) handler.TerminatableUpload {
	return upload.(*fileUpload)
}

func (store FileStore) AsLengthDeclarableUpload(upload handler.Upload) handler.LengthDeclarableUpload {
	return upload.(*fileUpload)
}

func (store FileStore) AsConcatableUpload(upload handler.Upload) handler.ConcatableUpload {
	return upload.(*fileUpload)
}

func (store FileStore) AsServableUpload(upload handler.Upload) handler.ServableUpload {
	return upload.(*fileUpload)
}

// defaultBinPath returns the path to the file storing the binary data
func (store FileStore) defaultBinPath(id string) string {
	return filepath.Join(store.Path, id)
}

type fileUpload struct {
	// info stores the current information about the upload
	info handler.FileInfo
	// redisKey is the Redis key for storing metadata
	redisKey string
	// binPath is the path to the binary file (which has no extension)
	binPath string
	// redisClient is the Redis client for metadata operations
	redisClient *redis.Client
}

func (upload *fileUpload) GetInfo(ctx context.Context) (handler.FileInfo, error) {
	return upload.info, nil
}

func (upload *fileUpload) WriteChunk(ctx context.Context, offset int64, src io.Reader) (int64, error) {
	file, err := os.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return 0, fmt.Errorf("failed to open binary file for writing: %w", err)
	}
	// Avoid the use of defer file.Close() here to ensure no errors are lost

	n, err := io.Copy(file, src)
	upload.info.Offset += n

	// Update metadata in Redis
	if updateErr := upload.writeInfo(ctx); updateErr != nil {
		_ = file.Close()
		return n, fmt.Errorf("failed to update metadata after chunk write: %w", updateErr)
	}

	if err != nil {
		_ = file.Close()
		return n, err
	}

	return n, file.Close()
}

func (upload *fileUpload) GetReader(ctx context.Context) (io.ReadCloser, error) {
	return os.Open(upload.binPath)
}

func (upload *fileUpload) Terminate(ctx context.Context) error {
	// Delete binary file
	if err := os.Remove(upload.binPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to remove binary file: %w", err)
	}

	// Delete metadata from Redis
	if err := upload.redisClient.Del(ctx, upload.redisKey).Err(); err != nil && err != redis.Nil {
		return fmt.Errorf("failed to remove metadata from Redis: %w", err)
	}

	return nil
}

func (upload *fileUpload) ConcatUploads(ctx context.Context, uploads []handler.Upload) (err error) {
	file, err := os.OpenFile(upload.binPath, os.O_WRONLY|os.O_APPEND, defaultFilePerm)
	if err != nil {
		return fmt.Errorf("failed to open file for concatenation: %w", err)
	}
	defer func() {
		// Ensure that close error is propagated, if it occurs.
		if cerr := file.Close(); err == nil {
			err = cerr
		}
	}()

	for _, partialUpload := range uploads {
		if err = partialUpload.(*fileUpload).appendTo(file); err != nil {
			return fmt.Errorf("failed to append upload: %w", err)
		}
	}

	// Update metadata after concatenation
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stats after concatenation: %w", err)
	}

	upload.info.Offset = stat.Size()
	if err := upload.writeInfo(ctx); err != nil {
		return fmt.Errorf("failed to update metadata after concatenation: %w", err)
	}

	return nil
}

func (upload *fileUpload) appendTo(file *os.File) error {
	src, err := os.Open(upload.binPath)
	if err != nil {
		return fmt.Errorf("failed to open source file for append: %w", err)
	}
	defer func(src *os.File) {
		_ = src.Close()
	}(src)

	if _, err = io.Copy(file, src); err != nil {
		return fmt.Errorf("failed to copy data during append: %w", err)
	}

	return nil
}

func (upload *fileUpload) DeclareLength(ctx context.Context, length int64) error {
	upload.info.Size = length
	upload.info.SizeIsDeferred = false
	return upload.writeInfo(ctx)
}

func (upload *fileUpload) FinishUpload(ctx context.Context) error {
	return upload.redisClient.Expire(ctx, upload.redisKey, DefaultTTL).Err()
}

func (upload *fileUpload) ServeContent(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, upload.binPath)
	return nil
}

// writeInfo stores the upload metadata in Redis
func (upload *fileUpload) writeInfo(ctx context.Context) error {
	data, err := json.Marshal(upload.info)
	if err != nil {
		return fmt.Errorf("failed to marshal upload info: %w", err)
	}

	err = upload.redisClient.Set(ctx, upload.redisKey, string(data), DefaultTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to store upload info in Redis: %w", err)
	}

	return nil
}

// createFile creates the file with the content. If the corresponding directory does not exist,
// it is created. If the file already exists, its content is removed.
func createFile(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, defaultFilePerm)
	if err != nil {
		if os.IsNotExist(err) {
			// Create directory if it doesn't exist
			if err = os.MkdirAll(filepath.Dir(path), defaultDirectoryPerm); err != nil {
				return fmt.Errorf("failed to create directory for %s: %w", path, err)
			}

			// Try creating the file again
			file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, defaultFilePerm)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", path, err)
			}
		} else {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if content != nil {
		if _, err = file.Write(content); err != nil {
			return fmt.Errorf("failed to write content to file %s: %w", path, err)
		}
	}

	return nil
}

// uid returns a unique id. These ids consist of 128 bits from a
// cryptographically strong pseudo-random generator and are like uuids, but
// without the dashes and significant bits.
//
// See: http://en.wikipedia.org/wiki/UUID#Random_UUID_probability_of_duplicates
func uid() string {
	id := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		// This is probably an appropriate way to handle errors from our source
		// for random bits.
		panic(err)
	}
	return hex.EncodeToString(id)
}
