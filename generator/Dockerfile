FROM golang:1.27-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
COPY internal/generator/ ./internal/generator/
COPY generator/main.go ./generator/main.go

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o /out/generator ./generator

FROM scratch
COPY --from=builder /out/generator /app/generator
ENTRYPOINT ["/app/generator"]
