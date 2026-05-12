# Backend

Данная папка содержит реализацию backend-части для программного продукта. Реализация выполнена при
помощи языка программирования Go, базы данных PostgreSQL (драйвер [pgx](https://github.com/jackc/pgx),
SQL-сборщик [squirrel](https://github.com/Masterminds/Squirrel)) и HTTP-фреймворка
[gin](https://github.com/gin-gonic/gin).

### Конфигурация

Для простоты настройки и лёгкой масштабируемости, в данном модуле используется конфигурация через
переменные окружения. Найти полный список и документацию к каждому из них можно в файле
[`.env.example`](.env.example). Все переменные должны иметь префикс `BACKEND_`.

### Начало работы

Данная документация подразумевает установленный Docker и поддержку плагина Docker Compose. Он
потребуется для запуска и сборки модулей проекта (клиентской и серверной частей) и запуска
необходимых сервисов, используемых проектом (PostgreSQL и MinIO).

Шаг 0. Склонируйте данный репозиторий, переместитесь в папку серверной части.
```bash
git clone https://github.com/Pelfox/edutrack.git
cd edutrack/backend/
```

Шаг 1. Установите зависимости проекта.
```bash
go mod download
```

Шаг 2. Настройте переменные окружения, скопировав и отредактировав файл с примерами.
```bash
cp .env.example .env
```

Шаг 3. Запустите сервер.
```bash
go run ./cmd/main.go
```

### OpenAPI

Документация API доступна после запуска сервера по адресу `/swagger/index.html`. Спецификация
доступна по адресу `/openapi.json` и
генерируется при помощи `swag`:
```bash
swag init --v3.1 --parseInternal --generalInfo main.go --dir cmd,internal/controllers,internal/dto,internal/repositories --output docs --outputTypes go,json,yaml
```

### Сборка через Docker

Для каждой версии, прошедшей интеграционное тестирование и готовой к выпуску, создаются готовые
Docker-образы, рекомендуется использовать именно их для запуска в production-среде:
```bash
docker pull ghcr.io/pelfox/edutrack/backend:latest
```

В случае необходимости локальной сборки, можно выполнить следующую команду:
```bash
docker build -f backend/Dockerfile -t edutrack/backend:latest backend/
```
