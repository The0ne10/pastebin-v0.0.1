FROM bitnami/minio:latest

# Установка mc (MinIO Client)
USER root
RUN install_packages curl
RUN curl -O https://dl.min.io/client/mc/release/linux-amd64/mc
RUN chmod +x mc && mv mc /usr/local/bin/mc

# Копирование скрипта инициализации в контейнер
COPY app/docker/minio/init.sh /usr/local/bin/init.sh
RUN chmod +x /usr/local/bin/init.sh

# Вернуться к обычному пользователю
USER 1001

# Команда запуска, включая скрипт инициализации
ENTRYPOINT ["sh", "-c", "/opt/bitnami/scripts/minio/run.sh & sh /usr/local/bin/init.sh"]
