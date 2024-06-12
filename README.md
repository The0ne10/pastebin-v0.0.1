# pastebin-v0.0.1

## Для запуска приложения

### нужно сделать файл исполняемым
* app/docker/minio/init.sh

### Через утилиту task
1. cd app/
2. go-task build-minio
3. go-task run

### Без task
1. docker-compose -f docker-compose.minio.yaml up --build -d
2. cd app/
3. go run ./cmd/pastebin/main.go