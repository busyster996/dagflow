package utility

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var ServiceName = "dagflow"

func FileOrPathExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func EnsureDirExist(name string) error {
	if !FileOrPathExist(name) {
		return os.MkdirAll(name, os.ModePerm)
	}
	return nil
}

func SliceToStrMap(s []string) map[string]string {
	m := make(map[string]string)
	for _, v := range s {
		slice := strings.Split(v, "=")
		switch {
		case len(slice) > 2:
			m[slice[0]] = strings.Join(slice[1:], "=")
		case len(slice) == 2:
			m[slice[0]] = slice[1]
		case len(slice) == 1:
			m[v] = ""
		}
	}
	return m
}

func MapToSlice(m map[string]string) []string {
	var s []string
	for k, v := range m {
		s = append(s, k+"="+v)
	}
	return s
}

func ClearDir(path string) {
	_ = os.RemoveAll(path)
	_ = EnsureDirExist(path)
}

func PathEscape(s string) string {
	s = filepath.Clean(strings.TrimPrefix(s, ".."))
	if s == ".." {
		return ""
	}
	if !strings.HasPrefix(s, "..") {
		return s
	}
	return PathEscape(s)
}

func CheckDuplicate[T comparable](slice []T) []T {
	var dup []T
	seen := make(map[T]struct{})
	for _, v := range slice {
		if _, exists := seen[v]; exists {
			dup = append(dup, v)
			continue
		}
		seen[v] = struct{}{}
	}
	if len(dup) == 0 {
		return nil
	}
	return dup
}

func RemoveDuplicate[T comparable](slice []T) []T {
	allKeys := make(map[T]bool)
	var list []T
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}

	return list
}

func JoinWithInvisibleChar(strs ...string) string {
	return strings.Join(strs, "\uFEFF")
}

func SplitByInvisibleChar(str string) []string {
	return strings.Split(str, "\uFEFF")
}

func ContainsInvisibleChar(str string) bool {
	return strings.Contains(str, "\uFEFF")
}

func MergerContext(parent context.Context, extras ...context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancelCause(parent)
	var stopFuncs []func() bool
	for _, extra := range extras {
		stop := context.AfterFunc(extra, func() {
			cancel(extra.Err())
		})
		stopFuncs = append(stopFuncs, stop)
	}
	return ctx, func() {
		cancel(nil)
		for _, stop := range stopFuncs {
			stop()
		}
	}
}

func MD5(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %s", dst, err)
	}

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, sourceFileStat.Mode())
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
