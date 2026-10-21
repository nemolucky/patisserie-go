# Patisserie API 🍰

RESTful бэкенд-сервис для интернет-магазина кондитерской продукции, написанный на **Go** с использованием принципов **чистой архитектуры (Clean Architecture)**.

Сервис предоставляет полный функционал для управления каталогом десертов, деревом категорий, загрузкой и оптимизацией медиафайлов в S3-совместимое хранилище (MinIO), а также реализует безопасную беспарольную аутентификацию по номеру телефона через SMS-код и JWT-токены.

---

## 📌 Содержание

- [Ключевые возможности](#-ключевые-возможности)
- [Стек технологий](#-стек-технологий)
- [Архитектура проекта](#-архитектура-проекта)
- [Структура репозитория](#-структура-репозитория)
- [Переменные окружения](#-переменные-окружения)
- [Запуск и разработка](#-запуск-и-разработка)
  - [Запуск через Docker Compose](#запуск-через-docker-compose)
  - [Локальный запуск](#локальный-запуск)
  - [Миграции базы данных](#миграции-базы-данных)
- [Спецификация REST API](#-спецификация-rest-api)
  - [Аутентификация (/api/v1/auth)](#1-аутентификация-apiv1auth)
  - [Пользователи (/api/v1/users)](#2-пользователи-apiv1users)
  - [Каталог и товары (/api/v1/products)](#3-каталог-и-товары-apiv1products)
  - [Категории (/api/v1/categories)](#4-категории-apiv1categories)
  - [Медиафайлы (/api/v1/images)](#5-медиафайлы-apiv1images)
- [Фоновые воркеры](#-фоновые-воркеры)

---

## ✨ Ключевые возможности

- **Каталог товаров и десертов**:
  - Создание и отображение карточек товаров с описанием, ценой и статусом активности.
  - Учет пищевой и энергетической ценности (КБЖУ: калории, белки, жиры, углеводы, вес).
  - Быстрый нечеткий поиск по каталогу с помощью расширения PostgreSQL `pg_trgm` (GIN-индекс).
  - Пагинация каталога (`limit`, `offset`).
- **Иерархические категории**:
  - Поддержка бесконечной вложенности (дерево категорий `parent_id`).
  - Получение как плоского списка, так и скомпонованного дерева категорий (`/tree`).
- **Медиа и S3-хранилище (MinIO)**:
  - Загрузка изображений через `multipart/form-data` с валидацией размера и типов.
  - Привязка нескольких изображений к товару с настраиваемым порядком сортировки (`sort_order`).
  - Фоновый воркер очистки (`ImageCleanupWorker`), регулярно удаляющий «осиротевшие» изображения из MinIO и БД.
- **Аутентификация и безопасность**:
  - Беспарольный вход по SMS-коду с временным хранением и верификацией в Redis.
  - Двухэтапная регистрация (проверка кода -> заполнение профиля: ФИО, Telegram).
  - JWT-аутентификация: короткоживущий `access_token` в теле ответа + долгоживущий `refresh_token` в защищенных `HttpOnly` Cookie.
  - Управление профилем: получение текущего пользователя (`/me`), смена телефона через SMS, удаление аккаунта.
- **Надежность и транзакции**:
  - Единый `TransactionManager` на базе `pgxpool`, прозрачно пробрасывающий `pgx.Tx` через `context.Context`.
  - Все SQL-запросы хранятся в отдельных `.sql` файлах и компилируются в бинарник через `//go:embed`.
  - Корректный `Graceful Shutdown` с ожиданием завершения HTTP-запросов и фоновых воркеров.
  - Структурированное логирование на базе встроенного `log/slog`.

---

## 🛠 Стек технологий

| Категория | Технология |
|---|---|
| **Язык** | Go (версия 1.26+) |
| **HTTP-роутер** | [chi](https://github.com/go-chi/chi) |
| **База данных** | PostgreSQL 17 |
| **Драйвер БД / Пул** | [pgx v5](https://github.com/jackc/pgx) (`pgxpool`) |
| **Миграции** | [golang-migrate](https://github.com/golang-migrate/migrate) |
| **S3 Object Storage** | [MinIO Go SDK](https://github.com/minio/minio-go) |
| **Кэш / TTL сессии** | [Redis v8](https://github.com/go-redis/redis) |
| **Аутентификация** | [golang-jwt v5](https://github.com/golang-jwt/jwt) |
| **Конфигурация** | [cleanenv](https://github.com/ilyakaznacheev/cleanenv) |
| **Контейнеризация** | Docker, Docker Compose |

---

## 🏗 Архитектура проекта

Проект спроектирован по принципам **Clean Architecture** (чистая архитектура) и слабой связанности компонентов:

```
[ HTTP Requests ]
        │
        ▼
[ Handlers / Router (internal/handler) ] ── (Middlewares: Auth, Logger, Recoverer)
        │
        ▼
[ UseCases / Business Logic (internal/usecase) ]
        │
   ┌────┴───────────────────────────┬────────────────────────────┐
   ▼                                ▼                            ▼
[ Postgres Repositories ]    [ MinIO Storage ]            [ Redis Cache ]
(internal/adapter/postgres)  (internal/adapter/minio)    (internal/adapter/redis)
   │
   ▼
[ Transaction Manager ] (pgx.Tx in Context)
```

- **Domain (`internal/domain`)**: Чистые доменные модели (`Product`, `Category`, `User`, `NutritionFacts`, `TokenPair`), ошибки и типы, независимые от внешних пакетов.
- **UseCase (`internal/usecase`)**: Сценарии использования и бизнес-правила (создание товаров с транзакциями, регистрация, ротация токенов, чистка медиа).
- **Adapters (`internal/adapter`)**: Реализация интерфейсов взаимодействия с внешними системами (PostgreSQL, MinIO, Redis, SMS-провайдер).
- **Handler (`internal/handler`)**: HTTP-контроллеры, валидация входных данных, маппинг DTO и формирование HTTP-ответов.
- **Worker (`internal/worker`)**: Асинхронные фоновые задачи, работающие по расписанию через `time.Ticker` и завершающиеся по сигналу контекста.

---

## 📂 Структура репозитория

```text
├── cmd/
│   └── main.go                  # Точка входа: сборка DI, старт воркеров, HTTP-сервера и Graceful Shutdown
├── internal/
│   ├── adapter/
│   │   ├── minio/               # Адаптер S3 MinIO (загрузка/удаление файлов, инициализация бакетов)
│   │   ├── postgres/            # Адаптер PostgreSQL, пул pgx, TransactionManager, SQL-репозитории
│   │   │   └── repositories/
│   │   │       ├── category/    # CRUD категорий и SQL-запросы
│   │   │       ├── image/       # Работа с таблицей images и поиск неиспользуемых файлов
│   │   │       ├── products/    # Работа с товарами, КБЖУ, связями и SQL-запросы в query/
│   │   │       └── users/       # Хранилище пользователей
│   │   ├── redis/               # Хранение SMS-кодов и статусов верификации телефонов
│   │   └── sms_provider/        # Провайдер отправки SMS (мок для локальной разработки)
│   ├── config/                  # Парсинг настроек приложения из переменных окружения
│   ├── ctxutil/                 # Утилиты работы со значениями в context.Context (userID)
│   ├── domain/                  # Сущности предметной области (Product, User, Category, Claims и т.д.)
│   ├── handler/                 # HTTP-обработчики (Auth, Product, Category, Image, User)
│   │   └── middleware/          # HTTP Middleware (AuthMiddleware для валидации JWT)
│   ├── sms_verifier/            # Сервис генерации и валидации одноразовых SMS-кодов
│   ├── token_provider/          # Генерация и верификация Access и Refresh JWT-токенов
│   ├── usecase/                 # Сценарии бизнес-логики (auth, category, image, product, user)
│   └── worker/                  # Фоновые воркеры (ImageCleanupWorker)
├── migrations/                  # SQL-миграции схемы БД (up/down)
├── pkg/
│   └── logger/                  # Настройка структурированного slog-логгера
├── .env.example                 # Пример переменных окружения
├── docker-compose.yml           # Декларация контейнеров: App, PostgreSQL, MinIO, Migrate
├── Dockerfile                   # Многоэтапная сборка легковесного контейнера приложения
└── Makefile                     # Команды для миграций и обслуживания
```

---

## ⚙️ Переменные окружения

Создайте файл `.env` на основе примера:

```bash
cp .env.example .env
```

| Переменная | Описание | Значение по умолчанию в `.env.example` |
|---|---|---|
| **HOST** | Хост для запуска HTTP-сервера | `0.0.0.0` |
| **PORT** | Порт приложения | `8080` |
| **DB_HOST** | Хост PostgreSQL | `patisserie-postgres` (или `localhost`) |
| **DB_PORT** | Порт PostgreSQL | `5432` |
| **DB_USER** | Пользователь базы данных | `postgres` |
| **DB_PASSWORD** | Пароль пользователя БД | `postgres` |
| **DB_NAME** | Имя базы данных | `patisserie_db` |
| **DB_SSL_MODE** | Режим SSL для PostgreSQL | `disable` |
| **MINIO_ENDPOINT** | Адрес подключения к MinIO | `patisserie-minio:9000` (или `localhost:9000`) |
| **ACCESS_KEY_ID** | Логин администратора MinIO | `admin` |
| **SECRET_ACCESS_KEY** | Пароль администратора MinIO | `admin123` |
| **USE_SSL** | Использование SSL для S3 | `false` |
| **BUCKET_NAME** | Имя бакета для изображений | `patisserie-images` |
| **BUCKET_POLICY** | Политика доступа к бакету | `public` |
| **REDIS_ADDR** | Адрес подключения к Redis | `patisserie-redis:6379` (или `localhost:6379`) |
| **REDIS_EXPIRATION** | Время жизни SMS-кодов в Redis | `5m` |
| **ACCESS_SECRET** | Секретный ключ для подписи Access Token | `change_me_access_secret` |
| **ACCESS_EXPIRATION**| Время жизни Access Token | `15m` (или `3h`) |
| **REFRESH_SECRET** | Секретный ключ для подписи Refresh Token | `change_me_refresh_secret` |
| **REFRESH_EXPIRATION**| Время жизни Refresh Token | `7d` |

---

## 🚀 Запуск и разработка

### Запуск через Docker Compose

Для развертывания всего окружения в Docker:

1. Подготовьте файл `.env.docker` или скопируйте `.env.example`:
   ```bash
   cp .env.example .env.docker
   ```
2. Запустите стек:
   ```bash
   docker compose up --build
   ```

Сервисы будут доступны по следующим адресам:
- **API сервис**: `http://localhost:8080/api/v1`
- **PostgreSQL**: `localhost:5432`
- **MinIO API**: `http://localhost:9000`
- **MinIO Web Console**: `http://localhost:9001` (логин/пароль: `admin` / `admin123`)

> **Примечание**: Перед первым запуском убедитесь, что в `docker-compose.yml` указаны подходящие локальные директории для монтирования volumes (по умолчанию указаны `/home/nemolucky/docker/...`).

---

### Локальный запуск

1. **Запустите необходимые зависимости** (PostgreSQL, Redis, MinIO) в Docker:
   ```bash
   docker run -d --name local-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=patisserie_db postgres:17-alpine
   docker run -d --name local-redis -p 6379:6379 redis:7-alpine
   docker run -d --name local-minio -p 9000:9000 -p 9001:9001 -e MINIO_ROOT_USER=admin -e MINIO_ROOT_PASSWORD=admin123 minio/minio server /data --console-address ":9001"
   ```

2. **Примените миграции**:
   ```bash
   make migration-up
   ```

3. **Запустите приложение**:
   ```bash
   go run ./cmd/main.go
   ```

---

### Миграции базы данных

Управление миграциями осуществляется через CLI-утилиту [golang-migrate](https://github.com/golang-migrate/migrate) или Makefile:

```bash
# Применить все новые миграции
make migration-up

# Откатить последнюю миграцию
make migration-down

# Создать новую пару SQL-файлов миграции
make migration-create name=add_orders_table
```

---

## 📡 Спецификация REST API

Базовый путь для всех эндпоинтов: `/api/v1`

### 1. Аутентификация (`/api/v1/auth`)

#### Запрос SMS-кода
```http
POST /api/v1/auth/sms
Content-Type: application/json

{
  "phone": "+79991234567"
}
```
*Ответ:* `200 OK` (SMS-код выводится в лог/консоль через мок-провайдер).

---

#### Подтверждение SMS-кода
```http
POST /api/v1/auth/verify
Content-Type: application/json

{
  "phone": "+79991234567",
  "code": "1234"
}
```
*Ответ (если пользователь уже зарегистрирован):* `200 OK`
```json
{
  "registration_required": false,
  "access_token": "eyJhbGciOiJIUzI1NiIs..."
}
```
*(Также устанавливается защищенный HttpOnly Cookie `refresh_token`).*

*Ответ (если пользователь новый):* `200 OK`
```json
true
```
*(Требуется вызвать `/auth/register`).*

---

#### Завершение регистрации нового пользователя
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "full_name": "Иванов Иван Иванович",
  "phone": "+79991234567",
  "telegram": "@ivanov"
}
```
*Ответ:* `201 Created`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

#### Обновление пары токенов (Refresh)
```http
POST /api/v1/auth/refresh
Cookie: refresh_token=<токен>
```
*Ответ:* `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

---

#### Выход (Logout)
```http
POST /api/v1/auth/logout
```
*Ответ:* `204 No Content` (Cookie `refresh_token` инвалидируется).

---

### 2. Пользователи (`/api/v1/users`)
> Требуется заголовок авторизации: `Authorization: Bearer <access_token>`

- `GET /api/v1/users/me` — Получение профиля текущего пользователя.
- `POST /api/v1/users/me/phone/sms` — Запрос SMS-кода на новый номер телефона:
  ```json
  {
    "phone": "+79997654321"
  }
  ```
- `POST /api/v1/users/me/phone/verify` — Подтверждение кода и сохранение нового номера:
  ```json
  {
    "phone": "+79997654321",
    "code": "1234"
  }
  ```
- `DELETE /api/v1/users/me` — Удаление своего аккаунта.

---

### 3. Каталог и товары (`/api/v1/products`)

#### Получение списка товаров каталога
```http
GET /api/v1/products?limit=20&offset=0&search=чизкейк
```
*Ответ:* `200 OK`
```json
{
  "products": [
    {
      "id": "e4f8d1c9-72b1-4c62-8924-a74582f3efc2",
      "title": "Нью-Йорк Чизкейк",
      "price": 450,
      "image": {
        "id": "2bc99403-d250-4d43-98fe-0e4277bfe4e9",
        "key": "2bc99403-d250-4d43-98fe-0e4277bfe4e9.jpg",
        "sort_order": 0
      },
      "categories": ["Десерты", "Чизкейки"]
    }
  ],
  "limit": 20,
  "offset": 0
}
```

---

#### Получение детальной информации о товаре
```http
GET /api/v1/products/{id}
```
*Ответ:* `200 OK`
```json
{
  "id": "e4f8d1c9-72b1-4c62-8924-a74582f3efc2",
  "title": "Нью-Йорк Чизкейк",
  "description": "Классический сливочный десерт на песочной основе",
  "price": 450,
  "is_active": true,
  "facts": {
    "calories": 320,
    "protein": 6.5,
    "fat": 22.0,
    "carbohydrates": 26.0,
    "weight": 150.0
  },
  "images": [
    {
      "id": "2bc99403-d250-4d43-98fe-0e4277bfe4e9",
      "key": "2bc99403-d250-4d43-98fe-0e4277bfe4e9.jpg",
      "sort_order": 0
    }
  ],
  "categories": [
    {
      "id": "f5e921d7-2eb1-4d14-87cf-45d28b03e05a",
      "title": "Чизкейки",
      "parent": null
    }
  ]
}
```

---

#### Создание товара
```http
POST /api/v1/products
Content-Type: application/json

{
  "title": "Макарон Фисташка",
  "description": "Французское пирожное с фисташковым ганашем",
  "price": 180,
  "images": [
    "2bc99403-d250-4d43-98fe-0e4277bfe4e9"
  ],
  "categories": [
    "f5e921d7-2eb1-4d14-87cf-45d28b03e05a"
  ],
  "facts": {
    "calories": 140,
    "protein": 2.5,
    "fat": 7.0,
    "carbohydrates": 16.5,
    "weight": 40.0
  }
}
```
*Ответ:* `201 Created` (возвращает UUID созданного товара).

---

#### Включение / выключение активности товара
```http
PATCH /api/v1/products/active/{id}
```
*Ответ:* `200 OK` (возвращает новый статус активности: `true` / `false`).

---

### 4. Категории (`/api/v1/categories`)

- `GET /api/v1/categories` — Плоский список всех категорий.
- `GET /api/v1/categories/tree` — Древовидная иерархическая структура категорий:
  ```json
  [
    {
      "id": "root-uuid",
      "title": "Выпечка",
      "children": [
        {
          "id": "child-uuid",
          "title": "Круассаны",
          "children": []
        }
      ]
    }
  ]
  ```
- `POST /api/v1/categories` — Создать категорию:
  ```json
  {
    "title": "Эклеры",
    "parent": "00000000-0000-0000-0000-000000000000"
  }
  ```
- `PATCH /api/v1/categories/{id}` — Изменить название категории (`{"title": "Новое название"}`).
- `DELETE /api/v1/categories/{id}` — Удалить категорию.

---

### 5. Медиафайлы (`/api/v1/images`)

#### Загрузка изображений
```http
POST /api/v1/images
Content-Type: multipart/form-data; boundary=---boundary

---boundary
Content-Disposition: form-data; name="images"; filename="croissant.jpg"
Content-Type: image/jpeg

<binary data>
---boundary--
```
*Ограничения:* Максимальный размер файла — 10 МБ.
*Ответ:* `200 OK` — массив UUID сохраненных изображений, готовых к привязке к товару:
```json
[
  "2bc99403-d250-4d43-98fe-0e4277bfe4e9"
]
```

---

## 🔄 Фоновые воркеры

В приложении работает фоновый процесс:

- **`ImageCleanupWorker`**:
  - Запускается в отдельной горутине при старте приложения.
  - С периодичностью **1 раз в час** выполняет поиск изображений в БД, которые были загружены, но не были привязаны ни к одному товару в таблице `product_images`.
  - Удаляет найденные неиспользуемые объекты из бакета MinIO и удаляет соответствующие записи из PostgreSQL.
  - Корректно прерывается при остановке приложения через контекст `context.Context`.
