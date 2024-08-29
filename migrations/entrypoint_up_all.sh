#!/bin/bash

# Применим миграции
# ВАЖНО! Для локального запуска вызвать из корня проекта, предварительно подгрузив переменные окружения в терминал командой  `set -a; source <путь к .env файлу>; set +a`
# Установка goose: go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir ./migrations postgres "user=$AVATAR_SERVICE_POSTGRES_USER password=$AVATAR_SERVICE_POSTGRES_PASSWORD dbname=$AVATAR_SERVICE_POSTGRES_DB host=$AVATAR_SERVICE_POSTGRES_HOST port=$AVATAR_SERVICE_POSTGRES_PORT sslmode=disable" up