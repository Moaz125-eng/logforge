FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/logforge-server ./cmd/server
RUN CGO_ENABLED=0 go build -o /out/logforge-agent ./cmd/agent

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/logforge-server /usr/local/bin/logforge-server
COPY --from=build /out/logforge-agent /usr/local/bin/logforge-agent
EXPOSE 8080 9090
ENTRYPOINT ["logforge-server"]
