BEGIN;
-- Создание таблиц
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    user_role VARCHAR(20) NOT NULL CHECK (user_role IN ('admin', 'manager', 'viewer')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    item_name VARCHAR(255) NOT NULL,
    item_description TEXT,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS items_history (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES items(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    operation VARCHAR(20) NOT NULL CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
    old_value TEXT,
    new_value TEXT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_items_history_item_id ON items_history(item_id);
CREATE INDEX IF NOT EXISTS idx_items_history_user_id ON items_history(user_id);
CREATE INDEX IF NOT EXISTS idx_items_history_changed_at ON items_history(changed_at);

-- Триггеры
-- TODO: Add triggers

-- Пользователи
-- TODO: Add test admin, manager, viewer users

-- Права доступа
GRANT ALL PRIVILIGES ON ALL TABLES IN SCHEMA public TO warehouse_control_user;
GRANT ALL PRIVILIGES ON ALL SEQUENCES IN SCHEMA public TO warehouse_control_user;

COMMIT;