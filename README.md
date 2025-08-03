# LogForge

Distributed log aggregation and search platform written in Go. Agents stream logs to a central server where entries are parsed, indexed, compressed, and queried through HTTP APIs and live tail streams.

## Features

- TCP and HTTP log ingestion with concurrent handlers
- Structured parsing for JSON and plaintext logs
- Goroutine worker pipeline with buffered channels
- Inverted and timestamp search indexes
- Query API with filters, pagination, and time ranges
- Chunked gzip storage with retention policies
- WebSocket live tail and multi-node forwarding
- Prometheus metrics and load benchmarks

## Quick start

```bash
cp .env.example .env
make build
make run
curl http://localhost:8080/health
```

## Docker

```bash
docker compose up --build
```

## Layout

```
cmd/server      central ingestion and query API
cmd/agent       forwarding agent for peer nodes
cmd/bench       ingestion load testing harness
internal/       ingestion, parsing, pipeline, index, query, storage
pkg/logentry    shared log entry types
```

## Operations

### Health and metrics

```bash
curl http://localhost:8080/health
curl http://localhost:8080/metrics
curl http://localhost:8080/metrics/prometheus
```

### Ingestion

```bash
curl -X POST http://localhost:8080/ingest \
  -H 'Content-Type: application/json' \
  -H 'X-Log-Service: billing' \
  -d '{"level":"error","service":"billing","message":"payment declined"}'
```

### Query and live tail

```bash
curl 'http://localhost:8080/query?q=payment&page=1&size=25'
websocat 'ws://localhost:8080/stream/tail?service=billing'
```

### Benchmark

```bash
go run ./cmd/bench -target http://localhost:8080 -workers 16 -total 5000
go run ./cmd/bench -target http://localhost:8080 -profile -out bench-report
```

### Production deploy

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

Environment variables are documented in `.env.example`. Data is stored under `LOGFORGE_DATA_DIR` with gzip archives and retention sweeps enabled by default.

## API overview

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Node health check |
| `POST /ingest` | HTTP log ingestion |
| `GET /query` | Search with filters and pagination |
| `GET /stream/tail` | WebSocket live tail |
| `GET /metrics` | JSON throughput snapshot |
| `GET /nodes` | Registered forwarding nodes |

## License

MIT
