#!/bin/sh

sleep 10

# Настройка mc (MinIO Client)
mc alias set myminio http://localhost:9000 $MINIO_ROOT_USER $MINIO_ROOT_PASSWORD

# Создание пользователя с заданными ключами доступа
mc admin user add myminio $MINIO_ACCESS_KEY $MINIO_SECRET_KEY

# Создание бакета
mc mb myminio/$MINIO_BUCKET_NAME

# Назначение политик доступа к бакету
mc admin policy set myminio readwrite user=$MINIO_ACCESS_KEY

# Предоставление прав пользователю на бакет
mc admin policy attach myminio readwrite --user=$MINIO_ACCESS_KEY

# Сделать бакет публичным
mc policy set public myminio/$MINIO_BUCKET_NAME
