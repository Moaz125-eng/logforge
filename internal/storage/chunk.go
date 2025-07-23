package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Chunk struct {
	Path      string
	CreatedAt time.Time
	Count     int
}

type ChunkWriter struct {
	dir     string
	maxSize int
	current *os.File
	chunk   Chunk
}

func NewChunkWriter(dir string, maxEntries int) (*ChunkWriter, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &ChunkWriter{dir: dir, maxSize: maxEntries}, nil
}

func (w *ChunkWriter) Write(entry logentry.Entry) error {
	if w.current == nil || w.chunk.Count >= w.maxSize {
		if err := w.rotate(); err != nil {
			return err
		}
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if _, err := w.current.Write(data); err != nil {
		return err
	}
	w.chunk.Count++
	return nil
}

func (w *ChunkWriter) rotate() error {
	if w.current != nil {
		_ = w.current.Close()
	}
	name := filepath.Join(w.dir, time.Now().UTC().Format("20060102-150405")+".chunk")
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	w.current = f
	w.chunk = Chunk{Path: name, CreatedAt: time.Now().UTC(), Count: 0}
	return nil
}

func (w *ChunkWriter) Close() error {
	if w.current == nil {
		return nil
	}
	return w.current.Close()
}
