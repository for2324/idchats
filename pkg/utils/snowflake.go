package utils

import (
	"fmt"
	"time"
)

const (
	// 时间戳占用位数
	timeBits uint8 = 41
	// 数据中心占用位数
	datacenterBits uint8 = 5
	// 机器标识占用位数
	machineBits uint8 = 5
	// 序列号占用位数
	sequenceBits uint8 = 12
	// 最大数据中心ID
	maxDatacenterID int64 = -1 ^ (-1 << datacenterBits)
	// 最大机器标识ID
	maxMachineID int64 = -1 ^ (-1 << machineBits)
	// 序列号掩码
	sequenceMask int64 = -1 ^ (-1 << sequenceBits)
)

// 雪花算法生成器结构体
type Snowflake struct {
	datacenterID int64 // 数据中心ID
	machineID    int64 // 机器标识ID
	sequence     int64 // 序列号
	lastTime     int64 // 上次生成ID的时间戳
}

// 生成ID方法
func (s *Snowflake) NextID() int64 {
	now := time.Now().UnixNano() / 1e6
	if now < s.lastTime {
		panic("Clock moved backwards!")
	}

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			now = s.waitNextMilli()
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = now
	return (now << (datacenterBits + machineBits + sequenceBits)) |
		(s.datacenterID << (machineBits + sequenceBits)) |
		(s.machineID << sequenceBits) |
		s.sequence
}

// 等待下一毫秒
func (s *Snowflake) waitNextMilli() int64 {
	now := time.Now().UnixNano() / 1e6
	for now <= s.lastTime {
		now = time.Now().UnixNano() / 1e6
	}
	return now
}

// 创建雪花算法生成器实例
func NewSnowflake(datacenterID, machineID int64) (*Snowflake, error) {
	if datacenterID > maxDatacenterID || datacenterID < 0 {
		return nil, fmt.Errorf("datacenter ID must be between 0 and %d", maxDatacenterID)
	}
	if machineID > maxMachineID || machineID < 0 {
		return nil, fmt.Errorf("machine ID must be between 0 and %d", maxMachineID)
	}
	return &Snowflake{
		datacenterID: datacenterID,
		machineID:    machineID,
		sequence:     0,
		lastTime:     time.Now().UnixNano() / 1e6,
	}, nil
}
