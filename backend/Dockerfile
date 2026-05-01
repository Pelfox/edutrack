# syntax=docker/dockerfile:1.7

# === Шаг сборки исходного кода проекта в выполняемый файл. ===
FROM golang:1.26.2-alpine3.23 AS build

WORKDIR /src

# Установка корневых сертификатов для поддержки TLS и временных зон.
RUN apk add --no-cache ca-certificates tzdata

# Копирование и установка всех зависимостей проекта.
COPY go.mod go.sum ./
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    go mod download

# Сборка финального выполняемого файла проекта.
COPY . .
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
      -trimpath \
      -buildvcs=false \
      -ldflags="-s -w -buildid=" \
      -o /out/app \
      ./cmd


# === Шаг сборки выполняемого файла для проведения миграций. ===
FROM --platform=$BUILDPLATFORM golang:1.26.2-alpine3.23 AS migrate-build

ARG MIGRATE_VERSION=v4.19.1

RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go install \
      -trimpath \
      -ldflags="-s -w -buildid=" \
      -tags "pgx5 file" \
      github.com/golang-migrate/migrate/v4/cmd/migrate@${MIGRATE_VERSION}

# === Финальный исполняемый (runtime) шаг. ===
FROM scratch

# Копирование статических данных: корневых сертификатов и временных зон.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

# Копирование выполняемых файлов и миграций.
COPY --from=build /out/app /app
COPY --from=migrate-build /go/bin/migrate /migrate
COPY migrations /migrations

# Запуск как не root-пользователь.
USER 65532:65532

# Публикация стандартного порта сервера.
EXPOSE 8000

ENTRYPOINT ["/app"]
