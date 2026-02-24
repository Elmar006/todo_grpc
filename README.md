# Todo gRPC Service

Сервис управления задачами (TODO) на базе gRPC, написанный на Go с использованием SQLite в качестве базы данных.

## Описание

Проект представляет собой backend-сервис для управления задачами. Реализует полный CRUD (Create, Read, Update, Delete) через gRPC API.

## Технологии

- **Язык**: Go 1.25
- **RPC-фреймворк**: gRPC
- **База данных**: SQLite (modernc.org/sqlite - чистая реализация на Go)
- **Прото-буфер**: Protocol Buffers v3
- **Логирование**: Logrus
- **Контейнеризация**: Docker, Docker Compose
- **Генерация кода**: Buf

## Структура проекта

```
todo_grpc/
├── backend/
│   ├── cmd/
│   │   └── server/          # Точка входа приложения
│   ├── internal/
│   │   ├── config/          # Конфигурация приложения
│   │   ├── db/              # Инициализация базы данных
│   │   ├── grpc/
│   │   │   ├── handler/     # gRPC обработчики
│   │   │   └── interceptor/ # gRPC интерцепторы
│   │   ├── logger/          # Логирование
│   │   ├── model/           # Модели данных
│   │   ├── repository/      # Слой доступа к данным
│   │   └── service/         # Бизнес-логика
│   ├── proto/
│   │   ├── gen/             # Сгенерированный код
│   │   └── todoService/     # Proto-файлы
│   ├── .env                 # Переменные окружения
│   ├── buf.gen.yaml         # Конфигурация генерации Buf
│   ├── buf.yaml             # Конфигурация Buf
│   ├── docker-compose.yml   # Docker Compose конфигурация
│   ├── Dockerfile           # Docker образ
│   ├── go.mod               # Go модуль
│   └── server.exe           # Скомпилированный бинарник
└── docker-compose.yml       # Основная Docker Compose конфигурация
```

## Архитектура

Проект следует архитектуре с разделением на слои:

1. **Handler (grpc/handler)** - Обработка gRPC запросов, валидация, преобразование между proto и внутренними моделями
2. **Service (internal/service)** - Бизнес-логика приложения
3. **Repository (internal/repository)** - Доступ к базе данных
4. **Model (internal/model)** - Внутренние модели данных

## API

Сервис предоставляет следующие методы:

| Метод | Описание |
|-------|----------|
| CreateTask | Создание новой задачи |
| GetTask | Получение задачи по ID |
| ListTasks | Получение списка всех задач |
| UpdateTask | Обновление задачи |
| DeleteTask | Удаление задачи |

### Структура задачи (Task)

- `id` (int64) - Уникальный идентификатор
- `title` (string) - Заголовок задачи
- `description` (string) - Описание задачи
- `completed` (bool) - Статус выполнения
- `created_at` (string) - Дата создания (RFC3339)
- `updated_at` (string) - Дата обновления (RFC3339)

## Установка и запуск

### Требования

- Go 1.25 или выше
- Docker и Docker Compose (опционально)
- Buf (для генерации proto-кода)

### Запуск через Docker Compose

```bash
docker-compose up --build
```

Сервер будет доступен на порту 50051.

### Локальный запуск

1. Установите зависимости:

```bash
cd backend
go mod download
```

2. Запустите сервер:

```bash
go run ./cmd/server
```

3. Или скомпилируйте и запустите бинарный файл:

```bash
go build -o server ./cmd/server
./server
```

### Генерация Proto-кода

Для генерации Go-кода из proto-файлов используйте Buf:

```bash
cd backend
buf generate
```

## Конфигурация

Конфигурация осуществляется через переменные окружения:

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| GRPC_PORT | Порт gRPC сервера | 50051 |

Файл `.env` расположен в директории `backend/`.

## База данных

Данные хранятся в SQLite базе данных по пути `./data/todo.db`. При использовании Docker Compose данные сохраняются в volume `todo_data`.

### Схема таблицы task

```sql
CREATE TABLE task (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Тестирование

Запуск тестов:

```bash
cd backend
go test ./...
```

## Логирование

Проект использует библиотеку Logrus для логирования. Логи выводятся в текстовом формате с временными метками.

## Пример использования

Для взаимодействия с сервисом необходимо использовать gRPC клиент. Пример на Go:

```go
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := todoService.NewTodoServiceClient(conn)

// Создание задачи
resp, err := client.CreateTask(context.Background(), &todoService.CreateTaskRequest{
    Title:       "Новая задача",
    Description: "Описание задачи",
})
```

## Лицензия

Проект доступен без явной лицензии.

## Планы развития

В ближайших обновлениях планируется:

### Умный поиск и фильтрация

Реализация расширенной фильтрации для метода `ListTasks`:

- Фильтрация по статусу выполнения (completed/not completed)
- Фильтрация по дате создания (диапазон дат)
- Полнотекстовый поиск по заголовку и описанию
- Сортировка по различным полям (дата создания, заголовок, статус)
- Пагинация результатов (limit/offset)

Пример будущего API:

```protobuf
message ListTasksRequest {
    optional bool completed = 1;
    optional string date_from = 2;
    optional string date_to = 3;
    optional string search_query = 4;
    optional string sort_by = 5;
    optional int32 limit = 6;
    optional int32 offset = 7;
}
```

### Тестирование gRPC хендлеров

Добавление полного покрытия тестами для слоя handler:

- Unit-тесты для всех методов gRPC хендлеров
- Интеграционные тесты с использованием testcontainers
- Тесты обработки ошибок и граничных условий
- Тесты таймаутов и контекста
- Mock-объекты для зависимости service слоя

Используемые инструменты:

- `stretchr/testify` — фреймворк для тестирования
- `gomock` — генерация mock-объектов
- `testcontainers-go` — интеграционное тестирование
