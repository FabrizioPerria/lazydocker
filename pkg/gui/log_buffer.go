package gui

import (
	"strings"
	"sync"
)

type LogBuffer struct {
	lines []string
	mutx  sync.Mutex
}

func NewLogBuffer() *LogBuffer {
	return &LogBuffer{
		lines: make([]string, 0),
		mutx:  sync.Mutex{},
	}
}

func (b *LogBuffer) Write (p []byte) (n int, err error) {
	b.mutx.Lock()
	defer b.mutx.Unlock()

	chunks := strings.Split(strings.TrimRight(string(p), "\n"), "\n")
	b.lines = append(b.lines, chunks...)
	return len(p), nil
}

func (b *LogBuffer) GetLines() []string {
	b.mutx.Lock()
	defer b.mutx.Unlock()

	return append([]string(nil), b.lines...)
}

func (b *LogBuffer) Clear() {
	b.mutx.Lock()
	defer b.mutx.Unlock()

	b.lines = nil
}
