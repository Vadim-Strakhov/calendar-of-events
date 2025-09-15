package cmd

import (
	"strings"
	"sync"
	"time"
)

type Log struct {
	mu      sync.Mutex
	entries []string
}

func NewLog() *Log {
	return &Log{entries: make([]string, 0, 128)}
}

func (l *Log) Add(parts ...string) {
	l.mu.Lock()
	l.entries = append(l.entries, time.Now().Format("15:04:05")+" "+strings.Join(parts, ""))
	l.mu.Unlock()
}

func (l *Log) All() []string {
	l.mu.Lock()
	cp := make([]string, len(l.entries))
	copy(cp, l.entries)
	l.mu.Unlock()
	return cp
}
