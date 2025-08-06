package rotate

import (
	"io"
	"sync"
)

type ICipher interface {
	Encrypt([]byte) []byte
	Decrypt([]byte) []byte
	EncryptReader(io.Reader, io.Writer) error
	DecryptReader(io.Reader, io.Writer) error
}

// Cipher represents a Caesar cipher with rotation key
type cipher struct {
	key    []byte // 密钥
	keyLen int    // 密钥长度
	shifts []int  // 预计算的移位值（各密钥字节 mod 8）
}

func New(key []byte) ICipher {
	shifts := make([]int, len(key))
	for i, b := range key {
		shifts[i] = int(b) % 8
	}
	return &cipher{
		key:    key,
		keyLen: len(key),
		shifts: shifts,
	}
}

// EncryptReader 加密
func (c *cipher) EncryptReader(r io.Reader, w io.Writer) error {
	return c.stream(r, w, true)
}

// DecryptReader 解密
func (c *cipher) DecryptReader(r io.Reader, w io.Writer) error {
	return c.stream(r, w, false)
}

// Encrypt 加密
func (c *cipher) Encrypt(data []byte) []byte {
	return c.process(data, true)
}

// Decrypt 解密
func (c *cipher) Decrypt(data []byte) []byte {
	return c.process(data, false)
}

func (c *cipher) stream(r io.Reader, w io.Writer, encrypt bool) error {
	if c.key == nil {
		_, err := io.Copy(w, r)
		return err
	}
	buf := make([]byte, c.keyLen)
	if c.keyLen > 1024 {
		buf = make([]byte, c.keyLen*4)
	}

	var err error
	for {
		var n int
		n, err = r.Read(buf)
		if n > 0 {
			_, err = w.Write(c.process(buf[:n], encrypt))
			if err != nil {
				return err
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func (c *cipher) process(data []byte, encrypt bool) []byte {
	if c.key == nil {
		return data
	}

	dataLen := len(data)
	result := make([]byte, dataLen)
	var wg sync.WaitGroup

	numChunks := (dataLen + c.keyLen - 1) / c.keyLen
	for i := 0; i < numChunks; i++ {
		start := i * c.keyLen
		end := start + c.keyLen
		if end > dataLen {
			end = dataLen
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()

			if encrypt {
				c.encryptChunk(data, result, start, end)
			} else {
				c.decryptChunk(data, result, start, end)
			}
		}(start, end)
	}
	wg.Wait()
	return result
}

// encryptChunk 加密数据块
func (c *cipher) encryptChunk(data, result []byte, start, end int) {
	key, shifts, keyLen := c.key, c.shifts, c.keyLen
	if keyLen == 1 {
		// 特化处理单字节密钥
		keyByte := key[0]
		shift := shifts[0]
		for j := start; j < end; j++ {
			b := data[j] ^ keyByte
			result[j] = (b >> shift) | (b << (8 - shift))
		}
		return
	}

	keyIndex := start % keyLen
	for j := start; j < end; j++ {
		// 异或和循环右移
		b := data[j] ^ key[keyIndex]
		shift := shifts[keyIndex]
		result[j] = (b >> shift) | (b << (8 - shift))

		keyIndex++
		if keyIndex >= keyLen {
			keyIndex = 0
		}
	}
}

// decryptChunk 解密数据块
func (c *cipher) decryptChunk(data, result []byte, start, end int) {
	key, shifts, keyLen := c.key, c.shifts, c.keyLen
	if keyLen == 1 {
		// 特化处理单字节密钥
		keyByte := key[0]
		shift := shifts[0]
		for j := start; j < end; j++ {
			b := data[j]
			b = (b << shift) | (b >> (8 - shift))
			result[j] = b ^ keyByte
		}
		return
	}

	keyIndex := start % keyLen
	for j := start; j < end; j++ {
		// 循环左移和异或
		b := data[j]
		shift := shifts[keyIndex]
		b = (b << shift) | (b >> (8 - shift))
		result[j] = b ^ key[keyIndex]

		keyIndex++
		if keyIndex >= keyLen {
			keyIndex = 0
		}
	}
}
