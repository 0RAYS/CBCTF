FROM node:22-alpine AS frontend-builder
WORKDIR /app

RUN corepack enable

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack use pnpm@10.30.2 && pnpm install --frozen-lockfile

COPY frontend/ .
RUN pnpm run build


FROM golang:1.26-alpine AS backend-builder
WORKDIR /app

RUN apk add --no-cache build-base libpcap-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=frontend-builder /app/dist ./frontend/dist
RUN CGO_ENABLED=1 go build -ldflags="-linkmode external -extldflags '-static' -s -w" -trimpath -o CBCTF .


FROM alpine:latest
WORKDIR /app

COPY --from=backend-builder /app/CBCTF .
RUN chmod +x /app/CBCTF

EXPOSE 8000

CMD ["/app/CBCTF"]
