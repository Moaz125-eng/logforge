package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Moaz125-eng/logforge/internal/config"
	"github.com/Moaz125-eng/logforge/internal/forward"
	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func main() {
	cfg := config.Load()
	agent := forward.NewAgent(cfg.ForwardPeers)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			entry := logentry.New(cfg.NodeID, "heartbeat", logentry.LevelInfo)
			_ = agent.Forward(ctx, entry)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("agent stopped")
}
