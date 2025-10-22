package sid

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type ID uint64

const (
	// Epoch: 2024-11-01 00:00:00.000 UTC
	epochUnix     uint64 = 1730390400000
	timestampBits uint64 = 42
	kindBits      uint64 = 4
	nodeBits      uint64 = 6
	seqBits       uint64 = 11

	maxKindID   uint64 = 1<<kindBits - 1
	maxNodeID   uint64 = 1<<nodeBits - 1
	maxSequence uint64 = 1<<seqBits - 1

	timeShift        = seqBits + nodeBits + kindBits
	kindShift        = seqBits + nodeBits
	nodeShift        = seqBits
	maxID     uint64 = (1 << (timestampBits + kindBits + nodeBits + seqBits)) - 1
)

// Snowflake generates unique snowflake IDs.
type Snowflake struct {
	mu        sync.Mutex
	kindNode  uint64
	sequence  uint64
	lastStamp uint64
}

// NewSnowflake returns a new Snowflake for the given kindID (0..15) and nodeID (0..63).
func NewSnowflake(kindID, nodeID uint64) (*Snowflake, error) {
	if kindID < 0 || kindID > maxKindID {
		return nil, fmt.Errorf("kindID %d  out of range", kindID)
	}
	if nodeID < 0 || nodeID > maxNodeID {
		return nil, fmt.Errorf("nodeID %d out of range", nodeID)
	}
	s := &Snowflake{
		kindNode: uint64(kindID)<<kindShift | uint64(nodeID)<<nodeShift,
		sequence: 1,
	}
	s.lastStamp = s.currentMillis()
	return s, nil
}

// Next returns the next unique ID.
func (s *Snowflake) Next() ID {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查序列是否溢出
	if s.sequence > maxSequence {
		// 等待时间前进
		for s.lastStamp > s.currentMillis() {
			time.Sleep(time.Millisecond)
		}
		s.lastStamp++
		s.sequence = 1
	} else {
		s.sequence++
	}

	id := (s.lastStamp << timeShift) | s.kindNode | s.sequence
	return ID(id & maxID)
}

// currentMillis returns milliseconds since custom epoch
func (s *Snowflake) currentMillis() uint64 {
	return uint64(time.Now().UnixNano()/1e6) - epochUnix
}

// Helpers to decode ID

func (f ID) String() string { return strconv.FormatUint(uint64(f), 10) }
func (f ID) Uint64() uint64 { return uint64(f) }
func (f ID) Int64() int64   { return int64(f) }
func (f ID) Bytes() []byte  { return []byte(f.String()) }
func (f ID) Base64() string { return base64.StdEncoding.EncodeToString(f.Bytes()) }
func (f ID) Time() time.Time {
	secs := uint64(f) >> timeShift
	return time.Unix(int64(secs+epochUnix), 0)
}
func (f ID) Kind() int64     { return int64((uint64(f) >> kindShift) & maxKindID) }
func (f ID) Node() int64     { return int64((uint64(f) >> nodeShift) & maxNodeID) }
func (f ID) Sequence() int64 { return int64(uint64(f) & maxSequence) }

// Parse helpers
func ParseID(id uint64) ID              { return ID(id) }
func ParseString(id string) (ID, error) { v, err := strconv.ParseUint(id, 10, 64); return ID(v), err }
func ParseBytes(id []byte) (ID, error)  { return ParseString(string(id)) }
