package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr       string
	TCPAddr        string
	DataDir        string
	RetentionDays  int
	WorkerCount    int
	BatchSize      int
	BatchWindow    time.Duration
	ForwardPeers   []string
	NodeID         string
	ShutdownWait   time.Duration
}

func Load() Config {
	return Config{
		HTTPAddr:      envOr("LOGFORGE_HTTP_ADDR", ":8080"),
		TCPAddr:       envOr("LOGFORGE_TCP_ADDR", ":9090"),
		DataDir:       envOr("LOGFORGE_DATA_DIR", "./data"),
		RetentionDays: envIntOr("LOGFORGE_RETENTION_DAYS", 14),
		WorkerCount:   envIntOr("LOGFORGE_WORKERS", 8),
		BatchSize:     envIntOr("LOGFORGE_BATCH_SIZE", 256),
		BatchWindow:   time.Duration(envIntOr("LOGFORGE_BATCH_WINDOW_MS", 50)) * time.Millisecond,
		ForwardPeers:  splitCSV(envOr("LOGFORGE_FORWARD_PEERS", "")),
		NodeID:        envOr("LOGFORGE_NODE_ID", "node-local"),
		ShutdownWait:  10 * time.Second,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOr(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func splitCSV(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
