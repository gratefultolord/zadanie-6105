# Tender Management Service

Этот проект представляет собой сервис для управления тендерами с использованием Go и базы данных PostgreSQL.

## Требования

Перед запуском проекта убедитесь, что у вас установлены следующие зависимости:

- [Go 1.19+](https://golang.org/doc/install)
- [PostgreSQL](https://www.postgresql.org/download/)
- [Docker](https://docs.docker.com/get-docker/) (для контейнеризации и деплоя)

## Настройка проекта через клонирование

### 1. Клонирование репозитория

```bash
git clone https://github.com/gratefultolord/zadanie-6105.git
cd zadanie-6105

### 2. Настройка переменных окружения

SERVER_ADDRESS=0.0.0.0:8080
POSTGRES_CONN=postgres://<username>:<password>@<host>:<port>/<dbname>
POSTGRES_JDBC_URL=postgresql://<host>:<port>/<dbname>
POSTGRES_USERNAME=<username>
POSTGRES_PASSWORD=<password>
POSTGRES_HOST=<host>
POSTGRES_PORT=5432
POSTGRES_DATABASE=<dbname>

### 3. Установка зависимостей

```bash
go mod tidy

### 4. Запуск проекта

```bash
go run cmd/app/main.go

## Настройка проекта через Docker

### 1. Построение Docker-образа

```bash
docker build -t tender-service .

### 2. Запуск контейнера

```bash
docker run --env-file .env -p 8080:8080 tender-service


## Тестирование

### 1. Проверка доступности сервера
- **Эндпоинт:** GET /ping
- **Цель:** Убедиться, что сервер готов обрабатывать запросы.
- **Ожидаемый результат:** Статус код 200 и текст "ok".

```yaml
GET /api/ping

Response:

  200 OK

  Body: ok
```

### 2. Тестирование функциональности тендеров
#### Получение списка тендеров
- **Эндпоинт:** GET /tenders
- **Описание:** Возвращает список тендеров с возможностью фильтрации по типу услуг.
- **Ожидаемый результат:** Статус код 200 и корректный список тендеров.

```yaml
GET /api/tenders

Response:

  200 OK

  Body: [ {...}, {...}, ... ]
```

#### Создание нового тендера
- **Эндпоинт:** POST /tenders/new
- **Описание:** Создает новый тендер с заданными параметрами.
- **Ожидаемый результат:** Статус код 200 и данные созданного тендера.

```yaml
POST /api/tenders/new

Request Body:

  {

    "name": "Тендер 1",

    "description": "Описание тендера",

    "serviceType": "Construction",

    "status": "Open",

    "organizationId": 1,

    "creatorUsername": "user1"

  }

Response:

  200 OK

  Body: 
  
  { 
    "id": 1, 
    "name": "Тендер 1", 
    "description": "Описание тендера",
    ...
  }
```

#### Получение тендеров пользователя
- **Эндпоинт:** GET /tenders/my
- **Описание:** Возвращает список тендеров текущего пользователя.
- **Ожидаемый результат:** Статус код 200 и список тендеров пользователя.

```yaml
GET /api/tenders/my?username=user1

Response:

  200 OK

  Body: [ {...}, {...}, ... ]  
```

#### Редактирование тендера
- **Эндпоинт:** PATCH /tenders/{tenderId}/edit
- **Описание:** Изменение параметров существующего тендера.
- **Ожидаемый результат:** Статус код 200 и обновленные данные тендера.

```yaml
PATCH /api/tenders/1/edit

Request Body:

  {

    "name": "Обновленный Тендер 1",

    "description": "Обновленное описание"

  }

Response:

  200 OK

  Body: 
  { 
    "id": 1, 
    "name": "Обновленный Тендер 1", 
    "description": "Обновленное описание",
    ...
  }  
```

### 3. Тестирование функциональности предложений
#### Создание нового предложения
- **Эндпоинт:** POST /bids/new
- **Описание:** Создает новое предложение для существующего тендера.
- **Ожидаемый результат:** Статус код 200 и данные созданного предложения.

```yaml
POST /api/bids/new

Request Body:

  {

    "name": "Предложение 1",

    "description": "Описание предложения",

    "status": "Submitted",

    "tenderId": 1,

    "organizationId": 1,

    "creatorUsername": "user1"

  }

Response:

  200 OK

  Body: 
  { 
    "id": 1, 
    "name": "Предложение 1", 
    "description": "Описание предложения",
    ...
  }
```

#### Получение списка предложений пользователя
- **Эндпоинт:** GET /bids/my
- **Описание:** Возвращает список предложений текущего пользователя.
- **Ожидаемый результат:** Статус код 200 и список предложений пользователя.

```yaml
GET /api/bids/my?username=user1

Response:

  200 OK

  Body: [ {...}, {...}, ... ]
  ```
  
#### Получение списка предложений для тендера
- **Эндпоинт:** GET /bids/{tenderId}/list
- **Описание:** Возвращает предложения, связанные с указанным тендером.
- **Ожидаемый результат:** Статус код 200 и список предложений для тендера.

```yaml
GET /api/bids/1/list

Response:

  200 OK

  Body: [ {...}, {...}, ... ]
  ```
  
#### Редактирование предложения
- **Эндпоинт:** PATCH /bids/{bidId}/edit
- **Описание:** Редактирование существующего предложения.
- **Ожидаемый результат:** Статус код 200 и обновленные данные предложения.

```yaml
PATCH /api/bids/1/edit

Request Body:

  {

    "name": "Обновленное Предложение 1",

    "description": "Обновленное описание"

  }

Response:

  200 OK

  Body: 
  { 
    "id": 1, 
    "name": "Обновленное Предложение 1", 
    "description": "Обновленное описание",
    ...,
  }
```

#### Откат версии предложения
- **Эндпоинт:** PUT /bids/{bidId}/rollback/{version}
- **Описание:** Откатить параметры предложения к указанной версии.
- **Ожидаемый результат:** Статус код 200 и данные предложения на указанной версии.

```yaml
PUT /api/bids/1/rollback/2

Response:

  200 OK

  Body: 
  { 
    "id": 1, 
    "name": "Предложение 1 версия 2", 
    ...
  }
```

### 4. Тестирование функциональности отзывов
#### Просмотр отзывов на прошлые предложения
- **Эндпоинт:** GET /bids/{tenderId}/reviews
- **Описание:** Ответственный за организацию может посмотреть прошлые отзывы на предложения автора, который создал предложение для его тендера.
- **Ожидаемый результат:** Статус код 200 и список отзывов на предложения указанного автора.

```yaml
GET /api/bids/1/reviews?authorUsername=user2&organizationId=1

Response:

  200 OK

  Body: [ {...}, {...}, ... ]
```

