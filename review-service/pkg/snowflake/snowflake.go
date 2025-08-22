package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	workerIDBits     = uint64(5)  // 机器ID位数
	dataCenterIDBits = uint64(5)  // 数据中心ID位数
	sequenceBits     = uint64(12) // 序列号位数

	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits)
	maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits)
	sequenceMask    = int64(-1) ^ (int64(-1) << sequenceBits)

	workerIDShift      = sequenceBits
	dataCenterIDShift  = sequenceBits + workerIDBits
	timestampLeftShift = sequenceBits + workerIDBits + dataCenterIDBits

	// 起始时间戳 (可以自己定义，比如项目上线时间)
	twepoch = int64(1672531200000) // 2023-01-01 00:00:00
)

type Snowflake struct {
	mu           sync.Mutex
	lastStamp    int64
	workerID     int64
	dataCenterID int64
	sequence     int64
}

// NewSnowflake 创建一个新的Snowflake节点
func NewSnowflake(workerID, dataCenterID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("workerID out of range")
	}
	if dataCenterID < 0 || dataCenterID > maxDataCenterID {
		return nil, errors.New("dataCenterID out of range")
	}
	return &Snowflake{
		workerID:     workerID,
		dataCenterID: dataCenterID,
		lastStamp:    -1,
		sequence:     0,
	}, nil
}

// NextID 生成下一个唯一ID
func (s *Snowflake) NextID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1e6 // 毫秒
	if timestamp < s.lastStamp {
		// 处理时钟回拨
		timestamp = s.lastStamp
	}

	if s.lastStamp == timestamp {
		// 同一毫秒内，序列递增
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			for timestamp <= s.lastStamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 不同毫秒，重置序列
		s.sequence = 0
	}

	s.lastStamp = timestamp

	return ((timestamp - twepoch) << timestampLeftShift) |
		(s.dataCenterID << dataCenterIDShift) |
		(s.workerID << workerIDShift) |
		s.sequence
}
