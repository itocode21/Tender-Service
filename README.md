# Tender Service API

## Описание

Это API для управления организациями, тендерами и предложениями. API позволяет создавать, читать, обновлять и удалять данные организаций, тендеров и предложений, а так же устанавливать им разные статусы.

## Установка

1.  Установите Go (версия 1.21 или выше). Вы можете скачать Go с [https://go.dev/dl/](https://go.dev/dl/).
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
    DB_NAME=TENDER --> ВАЖНО ИСПОЛЬЗОВАТЬ ИМЯ ДБ  TENDER
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
    ```
5.  Для запуска приложения в dev режиме:
    ```bash
    make run
    ```  
## Тестирование API

Вы можете использовать Postman или cURL для тестирования API, отправляя запросы на `http://localhost:8080`.

### Примеры API запросов

#### Создание организации

**Метод:** `POST`

**URL:** `http://localhost:8080/organizations`

**Тело запроса (JSON):**

```json
{
  "name": "Test Organization",
  "description": "This is a test organization.",
  "type": "LLC"
}
```
### Пример ответа (JSON):
```json
{
    "id": 1,
    "name": "Test Organization",
    "description": "This is a test organization.",
    "type": "LLC",
    "created_at": "2024-01-19T21:07:35.766786Z",
    "updated_at": "2024-01-19T21:07:35.766786Z"
} 
```

#### Получение списка тендеров

**Метод:** `GET`

**URL:** `http://localhost:8080/tenders`

### Пример ответа (JSON):
```json
[
    {
        "id": 1,
        "name": "test tender",
        "description": "this is a test tender",
        "organization_id": 1,
        "publication_date": "2024-01-19T21:38:01.190237Z",
        "end_date": "2024-01-20T21:38:01.190237Z",
        "status": "published",
        "created_at": "2024-01-19T21:38:01.190237Z",
        "updated_at": "2024-01-19T21:38:01.190237Z"
    }
]
```

#### Получение списка предложений по тендеру


**Метод:** `GET`

**URL:** `http://localhost:8080/tenders/{tender_id}/proposals` (где `{tender_id}` - это id тендера)

### Пример ответа (JSON):
```json
[
    {
        "id": 1,
        "tender_id": 1,
        "organization_id": 1,
        "description": "proposal description",
        "price": 100,
        "status": "pending",
        "created_at": "2024-01-20T00:35:30.151671Z",
        "updated_at": "2024-01-20T00:35:30.151671Z"
    }
]
```

#### Создание тендера

**Метод:** `POST`

**URL:** `http://localhost:8080/tenders`

**Тело запроса (JSON):**

```json
{
    "name": "Test Tender",
    "description": "This is a test tender.",
    "organization_id": 1,
    "publication_date": "2024-01-20T12:00:00Z",
    "end_date": "2024-01-27T12:00:00Z"
}
```
### Пример ответа (JSON):
```json
{
    "id": 1,
    "name": "Test Tender",
    "description": "This is a test tender.",
    "organization_id": 1,
    "publication_date": "2024-01-20T12:00:00Z",
    "end_date": "2024-01-27T12:00:00Z",
    "status": "draft",
    "created_at": "2024-01-20T09:15:00Z",
    "updated_at": "2024-01-20T09:15:00Z"
}
```

#### Создание предложения

**Метод:** `POST`

**URL:** `http://localhost:8080/proposals`

**Тело запроса (JSON):**

```json
{
    "tender_id": 1,
    "organization_id": 1,
    "description": "proposal description",
    "price": 100
}
```
### Пример ответа (JSON):
```json
{
    "id": 1,
    "tender_id": 1,
    "organization_id": 1,
    "description": "proposal description",
    "price": 100,
    "status": "pending",
    "created_at": "2024-01-20T09:15:00Z",
    "updated_at": "2024-01-20T09:15:00Z"
}
```

#### Обновление организации

**Метод:** `PUT`

**URL:** `http://localhost:8080/organizations/{id}` где `{id}` - это id организации)

**Тело запроса (JSON):**

```json
{
  "name": "Updated Test Organization",
  "description": "This is an updated test organization.",
    "type": "JSC"
}
```
### Пример ответа (JSON):
```json
{
    "id": 1,
    "name": "Updated Test Organization",
    "description": "This is an updated test organization.",
    "type": "JSC",
    "created_at": "2024-01-19T21:07:35.766786Z",
    "updated_at": "2024-01-20T10:00:00Z"
}
```

#### Удаление организации

**Метод:** `DELETE`

**URL:** `http://localhost:8080/organizations/{id}` где `{id}` - это id организации)

**Тело запроса (JSON):**

### Пример ответа (JSON):
```json
{
 "message": "Organization deleted successfully"
}
```

#### Получение организации по ID

**Метод:** `GET`

**URL:** `http://localhost:8080/organizations/{id}` где `{id}` - это id организации)

### Пример ответа (JSON):
```json
{
    "id": 1,
    "name": "Test Organization",
    "description": "This is a test organization.",
    "type": "LLC",
    "created_at": "2024-01-19T21:07:35.766786Z",
    "updated_at": "2024-01-19T21:07:35.766786Z"
}
```

#### Обновление тендера

**Метод:** `PUT`

**URL:** `http://localhost:8080/tenders/{id}` (где `{id}` - это id тендера)

**Тело запроса (JSON):**

```json
{
    "name": "Updated Test Tender",
    "description": "This is updated test tender",
    "publication_date": "2024-01-22T12:00:00Z",
    "end_date": "2024-01-29T12:00:00Z"
}
```
### Пример ответа (JSON):
```json
{
    "id": 1,
    "name": "Updated Test Tender",
    "description": "This is updated test tender",
    "organization_id": 1,
    "publication_date": "2024-01-22T12:00:00Z",
    "end_date": "2024-01-29T12:00:00Z",
    "status": "draft",
    "created_at": "2024-01-20T09:15:00Z",
     "updated_at": "2024-01-20T10:00:00Z"
}
```

#### Обновление предложения

**Метод:** `PUT`

**URL:** `http://localhost:8080/proposals/{id}` (где `{id}` - это id предложения)

**Тело запроса (JSON):**

```json
{
    "description": "Updated proposal description",
    "price": 200
}
```
### Пример ответа (JSON):
```json
{
    "id": 1,
    "tender_id": 1,
    "organization_id": 1,
    "description": "Updated proposal description",
    "price": 200,
    "status": "pending",
    "created_at": "2024-01-20T09:15:00Z",
    "updated_at": "2024-01-20T10:00:00Z"
}
```
#### Больше информации по api и endpoint смотри в /api

### Конфигурация линтера

Для статического анализа кода используется `golangci-lint`. Конфигурационный файл находится в корне проекта и называется `.golangci.yml`.

```yaml
run:
  # Установите путь до директории проекта
  # dir: ./
  skip-dirs:
    - ./vendor

linters-settings:
  govet:
    check-shadowing: true

linters:
  enable:
    - govet
    - errcheck
    - staticcheck
    - gosimple
    - unused
    - ineffassign
  #  disable:
    #  - typecheck

issues:
  exclude:
    - 'file is not go source file' # ignore warning about generated files