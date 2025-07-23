package storage

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type Manager struct {
	mu      sync.Mutex
	chunk   *ChunkWriter
	gzip    *GzipStore
	dataDir string
	bytes   int64
}

func NewManager(dataDir string) (*Manager, error) {
	chunkDir := filepath.Join(dataDir, "chunks")
	gzipDir := filepath.Join(dataDir, "archive")
	chunk, err := NewChunkWriter(chunkDir, 512)
	if err != nil {
		return nil, err
	}
	gzip, err := NewGzipStore(gzipDir)
	if err != nil {
		return nil, err
	}
	return &Manager{chunk: chunk, gzip: gzip, dataDir: dataDir}, nil
}

func (m *Manager) Store(entry logentry.Entry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.chunk.Write(entry)
}

func (m *Manager) FlushBatch(entries []logentry.Entry) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	path, err := m.gzip.Persist(entries)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(path)
	if err == nil {
		m.bytes += info.Size()
	}
	return path, nil
}

func (m *Manager) BytesStored() int64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.bytes
}

func (m *Manager) Close() error {
	return m.chunk.Close()
}
