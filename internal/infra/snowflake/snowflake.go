package snowflake

import (
	"errors"
	"sync"
	"time"
)

const (
	workerIDBits  = 10
	sequenceBits  = 12
	maxWorkerID   = int64(-1) ^ (int64(-1) << workerIDBits)
	sequenceMask  = int64(-1) ^ (int64(-1) << sequenceBits)
	workerIDShift = sequenceBits
	timeShift     = sequenceBits + workerIDBits
	customEpochMs = int64(1704067200000) // 2024-01-01T00:00:00Z
)

type Generator struct {
	mu       sync.Mutex
	workerID int64
	lastMs   int64
	sequence int64
}

func New(workerID int64) (*Generator, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, errors.New("worker id out of range")
	}
	return &Generator{workerID: workerID}, nil
}

func (g *Generator) NextID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	nowMs := time.Now().UTC().UnixMilli()
	if nowMs == g.lastMs {
		g.sequence = (g.sequence + 1) & sequenceMask
		if g.sequence == 0 {
			for nowMs <= g.lastMs {
				nowMs = time.Now().UTC().UnixMilli()
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastMs = nowMs
	return ((nowMs - customEpochMs) << timeShift) | (g.workerID << workerIDShift) | g.sequence
}
