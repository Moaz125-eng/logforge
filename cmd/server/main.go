package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Moaz125-eng/logforge/internal/config"
	"github.com/Moaz125-eng/logforge/internal/forward"
	"github.com/Moaz125-eng/logforge/internal/metrics"
	"github.com/Moaz125-eng/logforge/internal/index"
	"github.com/Moaz125-eng/logforge/internal/ingest"
	"github.com/Moaz125-eng/logforge/internal/parser"
	"github.com/Moaz125-eng/logforge/internal/pipeline"
	"github.com/Moaz125-eng/logforge/internal/query"
	"github.com/Moaz125-eng/logforge/internal/server"
	"github.com/Moaz125-eng/logforge/internal/storage"
	"github.com/Moaz125-eng/logforge/internal/stream"
	"github.com/Moaz125-eng/logforge/pkg/logentry"
)

func main() {
	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parserSvc := parser.NewService()
	storageSvc, err := storage.NewService(cfg)
	if err != nil {
		log.Fatalf("storage init failed: %v", err)
	}
	streamSvc := stream.NewService()
	forwardSvc := forward.NewService(cfg)
	metricsSvc := metrics.NewService()
	indexSvc := index.NewService()
	pipeSvc := pipeline.NewService(cfg, parserSvc.Parse, func(e logentry.Entry) error {
		metricsSvc.Collector().IncIngested()
		indexSvc.Index(e)
		streamSvc.Publish(e)
		metricsSvc.Collector().SetStreams(int32(streamSvc.HubCount()))
		if err := forwardSvc.Agent().Forward(ctx, e); err == nil {
			metricsSvc.Collector().IncForwarded()
		}
		if err := storageSvc.Persist(e); err != nil {
			return err
		}
		metricsSvc.Collector().IncStored()
		return nil
	})
	pipeSvc.Start(ctx)
	innerSink := func(e logentry.Entry) error {
		pipeSvc.Process(e)
		return nil
	}
	pipelineSink := ingest.WrapSink(cfg, innerSink)
	pipelineSink.Start(ctx)
	ingestSvc := ingest.NewService(cfg, pipelineSink.Sink)
	if err := ingestSvc.Start(ctx); err != nil {
		log.Fatalf("tcp ingest failed: %v", err)
	}

	queryEngine := query.NewEngine(indexSvc.Store())
	mux := server.NewMux(cfg, ingestSvc, parserSvc, indexSvc, queryEngine, storageSvc, streamSvc, forwardSvc, metricsSvc)
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
	pipelineSink.Wait()
	pipeSvc.Close()
	_ = storageSvc.Close()
	_ = httpServer.Shutdown(shutdownCtx)
}
