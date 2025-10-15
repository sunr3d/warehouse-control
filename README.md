# Warehouse Control

**Warehouse Control** — это мини-система для управления складом с CRUD-операциями, историей изменений и ролевой моделью доступа. Приложение демонстрирует использование PostgreSQL триггеров для аудита изменений (антипаттерн) и реализует полный стек веб-приложения с JWT авторизацией.

## Функциональность

- **CRUD операции инвентаря**:

  - `POST /items` - создание товара
  - `GET /items` - получение списка товаров
  - `PUT /items/{id}` - обновление товара
  - `DELETE /items/{id}` - удаление товара

- **История изменений данных**: Автоматическое логирование всех изменений через PostgreSQL триггеры

  - Кто, когда, что изменил
  - Полные снимки данных (old_value, new_value)
  - API endpoint `GET /items/{id}/history`

- **Ролевая модель доступа**:

  - **Admin** - полный доступ (создание, редактирование, удаление, просмотр истории)
  - **Manager** - создание, редактирование, просмотр истории
  - **Viewer** - только просмотр товаров

- **JWT авторизация**: Роль передается в токене и проверяется при каждом запросе

- **Простой веб-интерфейс**:
  - Вход через выпадающий список пользователей
  - Таблица товаров с действиями по ролям
  - Модальное окно для просмотра истории изменений

## Структура проекта

```
├── cmd/app/                    # Точка входа
├── internal/
│   ├── config/                 # Конфигурация
│   ├── entrypoint/             # Сборка зависимостей
│   ├── handlers/               # HTTP обработчики
│   │   ├── middleware/         # JWT + RBAC middleware
│   │   └── models.go           # DTO модели
│   ├── infra/postgres/         # Репозитории PostgreSQL
│   ├── interfaces/             # Интерфейсы (services, infra)
│   ├── services/               # Бизнес-логика
│   └── server/                 # HTTP сервер
├── models/                     # Доменные модели
├── migrations/                 # SQL миграции
├── web/                        # Статические файлы
```

## Установка и запуск

### Быстрый старт

1. **Клонирование репозитория**:

   ```bash
   git clone https://github.com/sunr3d/warehouse-control.git
   cd warehouse-control
   ```

2. **Запуск с Docker**:

   ```bash
   make up
   ```

3. **Доступ к приложению**:
   - Веб-интерфейс: http://localhost:8080
   - PostgreSQL: localhost:5433

### Тестовые пользователи

| Username   | Password | Role    | Права доступа                     |
| ---------- | -------- | ------- | --------------------------------- |
| admin123   | password | admin   | Полный доступ                     |
| manager123 | password | manager | Создание, редактирование, история |
| viewer123  | password | viewer  | Только просмотр                   |

### Команды Make

```bash
make up          # Запуск приложения
make down        # Остановка приложения
make restart     # Перезапуск
make clean       # Полная очистка (включая данные БД)
make logs        # Просмотр логов
make test        # Запуск тестов
make mocks       # Генерация моков
```

## API Endpoints

### Публичные

- `POST /login` - авторизация пользователя

### Защищенные (требуют JWT токен)

#### Товары

- `GET /items` - список товаров (admin, manager, viewer)
- `POST /items` - создание товара (admin, manager)
- `PUT /items/{id}` - обновление товара (admin, manager)
- `DELETE /items/{id}` - удаление товара (admin)

#### История

- `GET /items/{id}/history` - история изменений товара (admin, manager)

## База данных

### Схема

```sql
-- Пользователи
users (id, username, password_hash, user_role)

-- Товары
items (id, item_name, item_description, quantity, created_at, updated_at)

-- История изменений (автоматически через триггеры)
items_history (id, item_id, user_id, operation, old_value, new_value, changed_at)
```

### Триггеры

Приложение использует PostgreSQL триггеры для автоматического логирования изменений:

```sql
-- Триггер для логирования изменений
CREATE TRIGGER item_history_trigger
AFTER INSERT OR UPDATE ON items
FOR EACH ROW EXECUTE FUNCTION log_item_changes();
```

**Примечание**: Использование триггеров для аудита является антипаттерном и используется только по запросу из ТЗ.

## Переменные окружения

```env
HTTP_PORT=8080
LOG_LEVEL=info
JWT_SECRET=your-secret-key
DB_DSN=postgres://user:pass@host:port/db?sslmode=disable
DB_MAX_OPEN_CONNS=10
DB_MAX_IDLE_CONNS=2
```

## Тестирование

Проект включает unit тесты для критической бизнес-логики:

```bash
# Запуск всех тестов
make test

# Запуск тестов с покрытием
go test -v -cover ./...

# Тесты конкретного пакета
go test -v ./internal/services/authsvc/
go test -v ./internal/services/inventorysvc/
```

### Покрытие тестами

- ✅ `authsvc` - авторизация и JWT
- ✅ `inventorysvc` - бизнес-логика управления товарами
- ✅ Моки для всех интерфейсов

## Особенности реализации

### Антипаттерн: PostgreSQL триггеры

Проект намеренно использует триггеры для демонстрации антипаттерна:

```sql
-- Функция логирования изменений
CREATE OR REPLACE FUNCTION log_item_changes()
RETURNS TRIGGER AS $$
DECLARE
    current_user_id INTEGER;
BEGIN
    current_user_id := COALESCE(current_setting('warehouse.user_id', TRUE)::INTEGER, 0);

    IF (TG_OP = 'INSERT') THEN
        INSERT INTO items_history (item_id, user_id, operation, old_value, new_value)
        VALUES (NEW.id, current_user_id, 'INSERT', NULL, row_to_json(NEW)::TEXT);
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO items_history (item_id, user_id, operation, old_value, new_value)
        VALUES (NEW.id, current_user_id, 'UPDATE', row_to_json(OLD)::TEXT, row_to_json(NEW)::TEXT);
    END IF;
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;
```

### Передача user_id в триггеры

Для корректной работы триггеров используется `SET`:

```go
// В репозитории
_, err = tx.ExecContext(ctx, fmt.Sprintf("SET warehouse.user_id = %d", userID))
```
