package utility

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/robfig/cron/v3"

	"github.com/busyster996/dagflow/pkg/logx"
	"github.com/busyster996/dagflow/pkg/osext"
	"github.com/busyster996/dagflow/pkg/selfupdate"
)

type selfUpdate struct {
	base     bool
	sHash    []byte
	sURL     string
	cron     *cron.Cron
	skipFunc func() bool
}

func StartSelfUpdate(uri string, skipFunc func() bool) {
	s := &selfUpdate{
		sURL:     strings.TrimSuffix(uri, "/"),
		cron:     cron.New(cron.WithSeconds()),
		skipFunc: skipFunc,
	}
	// 获取当前程序的hash
	s.localSha256sum()
	_, err := s.cron.AddFunc("@every 30m", func() {
		s.doUpdate()
	})
	if err != nil {
		logx.Warnln(err)
		return
	}
	s.cron.Start()
}

func (p *selfUpdate) doUpdate() {
	if p.sURL == "" {
		logx.Debugln("automatic updates are not turned on")
		return
	}

	if p.sHash == nil {
		p.localSha256sum()
	}

	// 获取远端hash, example: https://oss.yfdou.com/tools/dagflow/latest.sha256sum
	checksum, name := p.remoteSha256sum(fmt.Sprintf("%s/latest.sha256sum", p.sURL))
	if checksum == nil || p.sHash == nil || bytes.Equal(p.sHash, checksum) {
		// 获取不到本地, 远端hash, 或者两者相同不更新
		return
	}

	opts := selfupdate.Options{
		Checksum: checksum,
	}
	// 检查权限
	if err := opts.CheckPermissions(); err != nil {
		logx.Warnln(err)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s", p.sURL, name))
	if err != nil {
		logx.Warnln(err)
		return
	}
	if resp == nil {
		logx.Warnln("resp is nil")
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logx.Errorln(err)
		}
	}(resp.Body)

	err = selfupdate.PrepareAndCheckBinary(resp.Body, opts)
	if err != nil {
		logx.Warnln(err)
		return
	}

	if p.skipFunc() {
		return
	}

	// 更新文件
	err = selfupdate.CommitBinary(opts)
	if err != nil {
		logx.Errorln(err)
		return
	}

	// 依赖systemd或者service守护
	os.Exit(200)
}

func (p *selfUpdate) localSha256sum() {
	path, err := osext.Executable()
	if err != nil {
		logx.Warnln(err)
		return
	}
	p.base = strings.Contains(path, "base")
	hash := sha256.New()
	file, err := os.Open(path)
	if err != nil {
		logx.Warnln(err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logx.Errorln(err)
		}
	}(file)
	_, err = io.Copy(hash, file)
	if err != nil {
		logx.Warnln(err)
		return
	}
	p.sHash = hash.Sum([]byte{})
	return
}

func (p *selfUpdate) remoteSha256sum(url string) ([]byte, string) {
	resp, err := http.Get(url)
	if resp != nil {
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				logx.Errorln(err)
			}
		}(resp.Body)
	}
	if err != nil {
		logx.Warnln(err)
		return nil, ""
	}
	body := bufio.NewReader(resp.Body)
	for {
		line, _, err := body.ReadLine()
		if err == io.EOF {
			return nil, ""
		}
		if err != nil {
			logx.Warnln(err)
		}
		if bytes.Contains(line, []byte(runtime.GOOS)) &&
			bytes.Contains(line, []byte(runtime.GOARCH)) &&
			bytes.Contains(line, []byte("base")) == p.base {
			fields := bytes.Fields(line)
			if len(fields) != 2 {
				continue
			}
			src := fields[0]
			n, err := hex.Decode(src, src)
			if err == nil {
				return src[:n], filepath.Base(string(fields[1]))
			} else {
				logx.Warnln(err)
			}
		}
	}
}
