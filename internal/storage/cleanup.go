package storage

import (
	"os"
	"path/filepath"
	"sort"
)

type CleanupReport struct {
	Removed int
	Kept    int
}

func PruneEmptyDirs(root string) (CleanupReport, error) {
	report := CleanupReport{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() || path == root {
			return nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			if err := os.Remove(path); err != nil {
				return err
			}
			report.Removed++
		} else {
			report.Kept++
		}
		return nil
	})
	return report, err
}

func ListChunks(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0)
	for _, e := range entries {
		if !e.IsDir() {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}
