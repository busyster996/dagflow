package file

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/busyster996/dagflow/pkg/lockfile"
	"github.com/busyster996/dagflow/pkg/tus/locker"
)

type sLocker struct {
	path string
}

func New(path string) locker.ILocker {
	_ = os.MkdirAll(path, os.ModeDir)
	return &sLocker{
		path: path,
	}
}

func (s *sLocker) NewLock(id string) (locker.ILock, error) {
	path, err := filepath.Abs(filepath.Join(s.path, id+".lock"))
	if err != nil {
		return nil, err
	}
	// We use Lockfile directly instead of lockfile.New to bypass the unnecessary
	// check whether the provided path is absolute since we just resolved it
	// on our own.
	return &fileUploadLock{
		file:                 lockfile.Lockfile(path),
		requestReleaseFile:   filepath.Join(s.path, id+".stop"),
		holderPollInterval:   5 * time.Second,
		acquirerPollInterval: 2 * time.Second,
	}, nil
}

type fileUploadLock struct {
	file lockfile.Lockfile

	requestReleaseFile   string
	holderPollInterval   time.Duration
	acquirerPollInterval time.Duration
}

func (lock *fileUploadLock) Lock(ctx context.Context) error {
	for {
		err := lock.file.TryLock()
		if err == nil {
			// Lock has been acquired, so we are good to go.
			break
		}
		if errors.Is(err, lockfile.ErrNotExist) {
			// ErrNotExist means that the file was not visible on disk yet. This
			// might happen when the disk is under some load. Wait a short amount
			// and retry.
			select {
			case <-ctx.Done():
				return errors.New("lock request timed out")
			case <-time.After(10 * time.Millisecond):
				continue
			}
		}
		// If the upload ID uses a folder structure (e.g. projectA/upload1), the directory
		// (e.g. projectA) might not exist, if no such upload exists already. In those cases,
		// we just return ErrNotFound because no such upload exists anyways.
		// TODO: This assumes that filelocker is used mostly with filestore, which is likely
		// true, but does not have to be. If another storage backend is used, we cannot make
		// any assumption about the folder structure. As an alternative, we should consider
		// normalizing the upload ID to remove the folder structure as well an turn projectA/upload1
		// into projectA-upload1.
		if errors.Is(err, fs.ErrNotExist) {
			return errors.New("upload not found")
		}
		if !errors.Is(err, lockfile.ErrBusy) {
			// If we get something different than ErrBusy, bubble the error up.
			return err
		}

		_ = os.MkdirAll(filepath.Dir(lock.requestReleaseFile), os.ModeDir)
		// If we are here, the lock is already held by another entity.
		// We create the .stop file to signal the lock holder to release the lock.
		file, err := os.Create(lock.requestReleaseFile)
		if err != nil {
			return err
		}
		_ = file.Close()

		select {
		case <-ctx.Done():
			// Context expired, so we return a timeout
			return errors.New("lock request timed out")
		case <-time.After(lock.acquirerPollInterval):
			// Continue with the next attempt after a short delay
			continue
		}
	}

	return nil
}

func (lock *fileUploadLock) Unlock() {
	err := lock.file.Unlock()

	// A "no such file or directory" will be returned if no lockfile was found.
	// Since this means that the file has never been locked, we drop the error
	// and continue as if nothing happened.
	if os.IsNotExist(err) {
		err = nil
	}

	// Try removing the file that is used for requesting a release. The error is
	// ignored on purpose.
	_ = os.Remove(lock.requestReleaseFile)
}
