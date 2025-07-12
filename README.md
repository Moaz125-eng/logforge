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

## License

MIT
