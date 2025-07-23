package storage

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

type GzipStore struct {
	dir string
}

func NewGzipStore(dir string) (*GzipStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &GzipStore{dir: dir}, nil
}

func (s *GzipStore) Persist(entries []logentry.Entry) (string, error) {
	path := filepath.Join(s.dir, "batch.gz")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	gz := gzip.NewWriter(f)
	enc := json.NewEncoder(gz)
	for _, entry := range entries {
		if err := enc.Encode(entry); err != nil {
			_ = gz.Close()
			_ = f.Close()
			return "", err
		}
	}
	if err := gz.Close(); err != nil {
		_ = f.Close()
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return path, nil
}

func (s *GzipStore) Read(path string) ([]logentry.Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	dec := json.NewDecoder(gz)
	out := make([]logentry.Entry, 0)
	for {
		var entry logentry.Entry
		if err := dec.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return out, err
		}
		out = append(out, entry)
	}
	return out, nil
}

func (s *GzipStore) Dir() string {
	return s.dir
}
