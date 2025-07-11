package tus

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"strings"

	"github.com/tjfoc/gmsm/sm3"
	"golang.org/x/crypto/sha3"
)

type HashProvider struct {
	hasher hash.Hash
}

func (hp *HashProvider) Checksum() []byte {
	return hp.hasher.Sum(nil)
}

func (hp *HashProvider) ChecksumHEX() string {
	return hex.EncodeToString(hp.Checksum())
}

func (hp *HashProvider) ChecksumBase64() string {
	return base64.StdEncoding.EncodeToString(hp.Checksum())
}

func NewHash(algorithm string) (hash.Hash, error) {
	switch strings.ToLower(algorithm) {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha512":
		return sha512.New(), nil
	case "sha224":
		return sha256.New224(), nil
	case "sha384":
		return sha512.New384(), nil
	case "sha3-256":
		return sha3.New256(), nil
	case "sm3":
		return sm3.New(), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

// ShaSumWriter 流式写时计算文件校验和
type ShaSumWriter struct {
	io.Writer
	*HashProvider
}

func NewShaSumWriter(writer io.Writer, algorithm string) (*ShaSumWriter, error) {
	hasher, err := NewHash(algorithm)
	if err != nil {
		return nil, err
	}
	return &ShaSumWriter{
		Writer:       writer,
		HashProvider: &HashProvider{hasher},
	}, nil
}

func (ssw *ShaSumWriter) Write(p []byte) (int, error) {
	n, err := ssw.Writer.Write(p)
	if err != nil {
		return n, err
	}

	if _, err = ssw.hasher.Write(p[:n]); err != nil {
		return n, err
	}

	return n, nil
}

// ShaSumReader 流式读时计算文件校验和
type ShaSumReader struct {
	io.Reader
	*HashProvider
}

func NewShaSumReader(algorithm string, reader io.Reader) (*ShaSumReader, error) {
	hasher, err := NewHash(algorithm)
	if err != nil {
		return nil, err
	}
	return &ShaSumReader{
		Reader:       reader,
		HashProvider: &HashProvider{hasher},
	}, nil
}

func (ssr *ShaSumReader) Read(p []byte) (int, error) {
	n, err := ssr.Reader.Read(p)
	if n > 0 {
		if _, err = ssr.hasher.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return n, err
}
