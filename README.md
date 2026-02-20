<<<<<<< HEAD
# Text Analysis Microservices

Два микросервиса для анализа текста.

## Архитектура

**Сервис A (Text Receiver)** - порт 8080
- Принимает HTTP запросы с текстом от пользователей
- Генерирует уникальный ID для каждого запроса
- Отправляет текст в Сервис B для анализа
- Хранит статус обработки в памяти
- Возвращает результаты анализа по ID

**Сервис B (Text Analyzer)** - порт 8081
- Получает текст от Сервиса A
- Подсчитывает статистику: количество слов, символов, предложений, среднюю длину слова
- Возвращает результаты обратно в Сервис A

## API Endpoints

### Сервис A

- `POST /api/v1/text` - принять текст для анализа
  ```json
  Request: {"text": "Ваш текст здесь"}
  Response: {"id": "uuid-запроса"}
  ```

- `GET /api/v1/status/{id}` - проверить статус обработки
  ```json
  Response: {
    "id": "uuid",
    "status": "processing|completed|failed",
    "result": {
      "word_count": 5,
      "char_count": 25,
      "sentence_count": 2,
      "average_word_len": 5.0
    },
    "created_at": "...",
    "updated_at": "..."
  }
  ```

- `GET /api/v1/health` - проверка здоровья сервиса
  ```json
  Response: {"status": "ok"}
  ```

### Сервис B

- `POST /api/v1/analyze` - анализ текста (используется Сервисом A)
- `GET /api/v1/health` - проверка здоровья сервиса

## Как запустить

### Вариант 1: Локально (без Docker)

**Шаг 1:** Запустите Сервис B (в первом терминале)
```bash
go run ./cmd/service-b
```
Вы увидите: `Service B starting on :8081`

**Шаг 2:** Запустите Сервис A (во втором терминале)
```bash
go run ./cmd/service-a
```
Вы увидите: `Service A starting on :8080`

**Важно:** Сервис A должен запускаться после Сервиса B, так как он обращается к нему по адресу `http://localhost:8081`

### Вариант 2: Через Docker Compose

```bash
# Собрать и запустить оба сервиса
docker-compose up --build

# Или в фоновом режиме
docker-compose up -d --build

# Остановить
docker-compose down
```

Docker Compose автоматически настроит сеть между сервисами и healthcheck.

## Как использовать

### 1. Отправить текст на анализ

```bash
curl -X POST http://localhost:8080/api/v1/text \
  -H "Content-Type: application/json" \
  -d "{\"text\": \"Hello world. How are you?\"}"
```

**Ответ:**
```json
{"id": "3dbdf3cb-76a8-4926-9827-67c81e20f4e4"}
```

### 2. Проверить статус обработки

Подставьте полученный `id` в запрос:

```bash
curl http://localhost:8080/api/v1/status/3dbdf3cb-76a8-4926-9827-67c81e20f4e4
```

**Ответ (если обработка завершена):**
```json
{
  "id": "3dbdf3cb-76a8-4926-9827-67c81e20f4e4",
  "status": "completed",
  "result": {
    "word_count": 5,
    "char_count": 25,
    "sentence_count": 2,
    "average_word_len": 5.0
  },
  "created_at": "2026-02-19T14:00:36Z",
  "updated_at": "2026-02-19T14:00:36Z"
}
```

**Ответ (если обработка еще идет):**
```json
{
  "id": "3dbdf3cb-76a8-4926-9827-67c81e20f4e4",
  "status": "processing",
  "created_at": "2026-02-19T14:00:36Z",
  "updated_at": "2026-02-19T14:00:36Z"
}
```

### 3. Проверить здоровье сервисов

```bash
# Сервис A
curl http://localhost:8080/api/v1/health

# Сервис B
curl http://localhost:8081/api/v1/health
```

## Пример полного цикла работы

```bash
# 1. Отправляем текст
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/text \
  -H "Content-Type: application/json" \
  -d "{\"text\": \"Привет мир. Как дела?\"}")

# 2. Извлекаем ID (в PowerShell)
$ID = ($RESPONSE | ConvertFrom-Json).id

# 3. Ждем немного и проверяем статус
Start-Sleep -Seconds 1
curl http://localhost:8080/api/v1/status/$ID
```

## Переменные окружения

- `SERVICE_B_URL` - URL Сервиса B (по умолчанию: `http://localhost:8081`)
  ```bash
  SERVICE_B_URL=http://service-b:8081 go run ./cmd/service-a
  ```

## Структура проекта

```
text-analysis-test-task/
├── cmd/
│   ├── service-a/          # Точка входа Сервиса A
│   └── service-b/          # Точка входа Сервиса B
├── pkg/
│   ├── servicea/           # Логика Сервиса A
│   │   ├── handlers.go     # HTTP хендлеры
│   │   ├── storage.go      # Хранилище статусов (in-memory)
│   │   └── client.go       # HTTP клиент к Сервису B
│   └── serviceb/           # Логика Сервиса B
│       ├── handlers.go     # HTTP хендлеры
│       └── analyzer.go     # Логика анализа текста
├── docker-compose.yml      # Конфигурация Docker Compose
├── Dockerfile.service-a    # Dockerfile для Сервиса A
└── Dockerfile.service-b    # Dockerfile для Сервиса B
```

## Особенности реализации

- ✅ Асинхронная обработка (Сервис A не блокируется при отправке в Сервис B)
- ✅ In-memory хранилище с thread-safe доступом (`sync.RWMutex`)
- ✅ Graceful shutdown (корректное завершение при Ctrl+C)
- ✅ Таймауты HTTP запросов (10 секунд)
- ✅ Обработка ошибок сети
- ✅ Логирование операций
- ✅ Health check endpoints

=======
# text-analysis-test-task
>>>>>>> ad67d66ddeb9f0ec5cbae21b1cea30df15ca7348
