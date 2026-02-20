
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


## Запуск

### Локально 

**1** Запустить Сервис B
```bash
go run ./cmd/service-b
```

**2** Запустить Сервис A
```bash
go run ./cmd/service-a
```

### Через Docker Compose

```bash
docker-compose up --build
```


## Использование

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

**Ответ:**
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

### 3. Проверить здоровье сервисов

```bash
# Сервис A
curl http://localhost:8080/api/v1/health

# Сервис B
curl http://localhost:8081/api/v1/health
```

