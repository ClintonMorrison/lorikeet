# Build stage
FROM golang:1.25.6 AS builder

WORKDIR /build

COPY ./go.mod ./
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN go build -o server cmd/server/main.go
RUN go build -o doBackup cmd/backup/main.go
RUN go build -o doRestore cmd/restore/main.go

# Runtime stage
FROM ubuntu:24.04

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        nginx \
        ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /build/server ./
COPY --from=builder /build/doBackup ./
COPY --from=builder /build/doRestore ./
COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./ui/build ./ui/build
COPY ./start.sh ./

CMD ["sh", "./start.sh"]
