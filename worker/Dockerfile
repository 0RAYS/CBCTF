FROM golang:1.27-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
COPY internal/worker/ ./internal/worker/
COPY worker/main.go ./worker/main.go

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o /out/worker ./worker

FROM scratch
COPY --from=builder /out/worker /app/worker
ENTRYPOINT ["/app/worker"]
