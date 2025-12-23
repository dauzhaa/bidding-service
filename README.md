# Art Bidding Microservices System

Распределенная система для проведения онлайн-аукционов. Реализован полный цикл работы: от создания лота на основе данных из Metropolitan Museum of Art до обработки ставок с защитой от состояния гонки (Race Conditions) и асинхронной отправки уведомлений.

## Функциональные возможности

* **Микросервисная архитектура:** Система разделена на 4 сервиса (Bidding, Auction, Notification, Gateway).
* **Транспорт:** Взаимодействие между сервисами реализовано через gRPC с использованием Protobuf.
* **Асинхронность:** Обработка событий реализована через брокер сообщений Apache Kafka.
* **Внешние интеграции:** Взаимодействие с публичным API The Met Museum для получения данных о предметах искусства.
* **Работа под нагрузкой:** Использование распределенных блокировок (Redis Distributed Lock) для обеспечения консистентности данных при конкурентных запросах.
* **Шардирование БД:** Поддержка горизонтального масштабирования базы данных PostgreSQL.
* **API Gateway:** Трансляция REST (HTTP) запросов в gRPC, документация в формате Swagger.
* **Тестирование:** Unit-тесты с использованием Mocks и Suites (покрытие бизнес-логики >84%).

## Технологический стек

| Категория | Технологии |
| :--- | :--- |
| **Язык** | Go (Golang) 1.22+ |
| **Базы данных** | PostgreSQL (pgx driver), Redis |
| **Брокер сообщений** | Apache Kafka, Zookeeper |
| **Коммуникация** | gRPC, HTTP/1.1 (REST) |
| **Инфраструктура** | Docker, Docker Compose |
| **Инструменты** | Make, Mockery, Swagger (OpenAPI v2), Testify |

## Инструкция по запуску

### 1. Запуск системы
Вся инфраструктура и микросервисы запускаются в Docker-контейнерах.

```bash
docker compose up -d --build
```

2. Инициализация базы данных

Необходимо создать таблицы в PostgreSQL (выполняется один раз после первого запуска контейнеров).

Таблица для аукционов (Шард 1):

```Bash
docker exec -i goproject-postgres-shard-1-1 psql -U user -d auction_db_1 -c "CREATE TABLE IF NOT EXISTS auctions (id BIGINT PRIMARY KEY, title TEXT, artist TEXT, start_price BIGINT, image_url TEXT, status TEXT);"
```

Таблица для ставок (Шард 1 и Шард 2):

```Bash
docker exec -i goproject-postgres-shard-1-1 psql -U user -d auction_db_1 -c "CREATE TABLE IF NOT EXISTS bids (id SERIAL PRIMARY KEY, auction_id BIGINT, user_id BIGINT, amount BIGINT, created_at TIMESTAMP);"

docker exec -i goproject-postgres-shard-2-1 psql -U user -d auction_db_2 -c "CREATE TABLE IF NOT EXISTS bids (id SERIAL PRIMARY KEY, auction_id BIGINT, user_id BIGINT, amount BIGINT, created_at TIMESTAMP);"
```

Тестирование API
Создание аукциона

Создание лота на основе картины Клода Моне "Букет подсолнухов" (ID музея: 436121):

```Bash
curl -X POST http://localhost:8081/v1/auctions \
  -d '{"object_id": 436121, "start_price": 5000}'
```
Размещение ставки

Отправка ставки на созданный лот:

```Bash
curl -X POST http://localhost:8081/v1/bids \
  -d '{"auction_id": 436121, "user_id": 1, "amount": 6000}'
```

Просмотр уведомлений

Логи сервиса уведомлений можно просмотреть командой:

```Bash
docker logs -f goproject-notification-service-1
```

Документация API

Спецификация API генерируется автоматически из .proto файлов. Для просмотра:

    Откройте Swagger Editor.

    Импортируйте файл api/swagger/auction_service.swagger.json или bidding_service.swagger.json из репозитория проекта.