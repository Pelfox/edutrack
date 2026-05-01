# syntax=docker/dockerfile:1.7

# === Шаг сборки модуля ===
FROM node:lts-alpine3.23 AS build

WORKDIR /app

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

# Включаем поддержку пакетного менеджера pnpm.
RUN corepack enable

# Копируем файлы с зависимостями для улучшения кэширования.
COPY package.json pnpm-lock.yaml ./
COPY pnpm-workspace.yaml* ./
COPY .npmrc* ./

# Загружаем зависимости из скопированных файлов.
RUN --mount=type=cache,id=pnpm-store,target=/pnpm/store \
    pnpm config set store-dir /pnpm/store && \
    pnpm fetch --frozen-lockfile

COPY . .

# Устанавливаем ранее собранные и скачанные зависимости, собираем проект.
RUN --mount=type=cache,id=pnpm-store,target=/pnpm/store \
    pnpm install --offline --frozen-lockfile && \
    pnpm run build

# Компрессируем (через gzip) статические файлы сборки.
RUN find dist -type f \
    \( -name '*.js' -o -name '*.css' -o -name '*.html' -o -name '*.svg' -o -name '*.json' -o -name '*.txt' -o -name '*.xml' \) \
    -exec gzip -k -9 {} \;

# === Финальный шаг: минималистичный Nginx. ===
FROM nginxinc/nginx-unprivileged:stable-alpine AS runtime

# Копируем необходимые файлы для раздачи через Nginx.
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 8080

STOPSIGNAL SIGQUIT
