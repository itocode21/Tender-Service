# Tender Service API

## Описание

Это API для управления организациями, тендерами и предложениями. API позволяет создавать, читать, обновлять и удалять данные организаций, тендеров и предложений, а так же устанавливать им разные статусы.

## Установка

1.  Установите Go (версия 1.23 или выше). Вы можете скачать Go с [https://go.dev/dl/](https://go.dev/dl/).
2.  Установите Docker и Docker Compose. Docker Desktop можно скачать отсюда [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/).
3. Установите GolangCI-lint  с [https://github.com/golangci/golangci-lint#installation](https://github.com/golangci/golangci-lint#installation)
4.  Клонируйте репозиторий:

    ```bash
    git clone https://github.com/your-username/Tender-Service.git
    cd Tender-Service
    ```

## Настройка

1.  Создайте файл `.env` в корневой директории проекта и заполните переменные окружения:

    ```env
    PORT=8080
    DB_HOST=postgres
    DB_PORT=5432
    DB_USER=ito21
    DB_PASSWORD=1899
    DB_NAME=TENDER --> ВАЖНО ДБ ИМЯ ДОЛЖНО БЫТЬ ВСЕГДА TENDER, все сиквел запросы написанос этим условием
    ```
2. Убедитесь что значения в файле `docker-compose.yaml` совпадают с файлом `.env`

## Запуск

1.  Для запуска приложения с использованием Docker Compose:
    ```bash
    make up
    ```
2.  Для остановки приложения с использованием Docker Compose:
    ```bash
    make down
    ```
3. Для остановки и удаления контейнеров и volumes:
    ```bash
    make down-v
    ```
4. Для удаления docker сетей:
   ```bash
   make network-prune

   