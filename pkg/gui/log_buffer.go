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

	return append([]string{}, b.lines...)
}

func (b *LogBuffer) Update(linenumber int, line string) {
	b.mutx.Lock()
	defer b.mutx.Unlock()

	if linenumber >= 0 && linenumber < len(b.lines) {
		b.lines[linenumber] = line
	}
}

func (b *LogBuffer) HighlightKeywordInLine(linenumber int, keyword string) string {
	b.mutx.Lock()
	defer b.mutx.Unlock()
	line := b.lines[linenumber]
	if strings.Contains(line, keyword) {
		b.lines[linenumber] = strings.ReplaceAll(line, keyword, "\x1b[37;41m"+keyword+"\x1b[0m")
	}
	return b.lines[linenumber]
}

func (b *LogBuffer) ColorKeywordInLine(linenumber int, keyword string) string {
	b.mutx.Lock()
	defer b.mutx.Unlock()
	line := b.lines[linenumber]
	if strings.Contains(line, keyword) {
		b.lines[linenumber] = strings.ReplaceAll(line, keyword, "\x1b[31m"+keyword+"\x1b[0m")
	}
	return b.lines[linenumber]
}

func (b *LogBuffer) CleanLineFromColors(linenumber int) string {
	b.mutx.Lock()
	defer b.mutx.Unlock()

	if linenumber < 0 || linenumber >= len(b.lines) {
		return "" // out of bounds
	}
	line := b.lines[linenumber]
	// Remove ANSI color codes
	cleanedLine := strings.NewReplacer(
		"\x1b[37;41m", "",
		"\x1b[31m", "",
		"\x1b[0m", "",
	).Replace(line)
	return cleanedLine
}

func (b *LogBuffer) ColorToHighlightKeyword(linenumber int, keyword string) {
	b.mutx.Lock()
	defer b.mutx.Unlock()

	if linenumber >= 0 && linenumber < len(b.lines) {
		b.CleanLineFromColors(linenumber)
		b.HighlightKeywordInLine(linenumber, keyword)
	}
}

func (b *LogBuffer) HighlightToColorKeyword(linenumber int, keyword string) string {
	b.mutx.Lock()
	defer b.mutx.Unlock()

	if linenumber >= 0 && linenumber < len(b.lines) {
		b.CleanLineFromColors(linenumber)
		b.ColorKeywordInLine(linenumber, keyword)
	}
	return b.lines[linenumber]
}
	

func (b *LogBuffer) Clear() {
	b.mutx.Lock()
	defer b.mutx.Unlock()

	panic("Clear() is not implemented yet") // TODO: implement Clear method
	b.lines = nil
}
