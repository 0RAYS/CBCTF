FROM node:22-alpine AS frontend-builder
WORKDIR /app

RUN corepack enable

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack use pnpm@10.30.2 && pnpm install --frozen-lockfile

COPY frontend/ .
RUN pnpm run build


FROM golang:1.26-bookworm AS backend-builder
WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends libpcap-dev gcc

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=frontend-builder /app/dist ./frontend/dist

RUN go build -ldflags="-s -w" -trimpath -o CBCTF .


FROM alpine:3.21
WORKDIR /app

COPY --from=backend-builder /app/CBCTF .

EXPOSE 8000

CMD ["./CBCTF"]
