package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Moaz125-eng/logforge/internal/config"
	"github.com/Moaz125-eng/logforge/internal/ingest"
	"github.com/Moaz125-eng/logforge/internal/parser"
	"github.com/Moaz125-eng/logforge/internal/server"
	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func main() {
	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parserSvc := parser.NewService()
	sink := func(e logentry.Entry) error {
		_, err := parserSvc.Parse("json", e.Raw)
		return err
	}
	ingestSvc := ingest.NewService(cfg, sink)
	if err := ingestSvc.Start(ctx); err != nil {
		log.Fatalf("tcp ingest failed: %v", err)
	}

	mux := server.NewMux(cfg, ingestSvc, parserSvc)
	httpServer := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
	}

	go func() {
		log.Printf("logforge http listening on %s tcp on %s", cfg.HTTPAddr, cfg.TCPAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownWait)
	defer shutdownCancel()
	cancel()
	ingestSvc.Wait()
	_ = httpServer.Shutdown(shutdownCtx)
}
