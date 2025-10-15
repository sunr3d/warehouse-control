BEGIN;
-- Создание таблиц
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    user_role VARCHAR(20) NOT NULL CHECK (user_role IN ('admin', 'manager', 'viewer'))
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
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
CREATE OR REPLACE FUNCTION log_item_changes()
RETURNS TRIGGER AS $$
DECLARE
    current_user_id INTEGER;
BEGIN
    current_user_id := COALESCE(current_setting('warehouse.user_id', TRUE)::INTEGER, 0);

    IF (TG_OP = 'INSERT') THEN
        INSERT INTO items_history (item_id, user_id, operation, old_value, new_value, changed_at)
        VALUES (NEW.id, current_user_id, 'INSERT', NULL, row_to_json(NEW)::TEXT, CURRENT_TIMESTAMP);
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO items_history (item_id, user_id, operation, old_value, new_value, changed_at)
        VALUES (NEW.id, current_user_id, 'UPDATE', row_to_json(OLD)::TEXT, row_to_json(NEW)::TEXT, CURRENT_TIMESTAMP);
    END IF;
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER item_history_trigger
AFTER INSERT OR UPDATE ON items
FOR EACH ROW EXECUTE FUNCTION log_item_changes();

-- Пользователи
INSERT INTO users (username, password_hash, user_role) VALUES
('admin123', '$2a$08$6Y7QSNKu4oJ55nKNf4iBQueMiT7Hfg5.3p.UtNCB0EhX3chigXYau', 'admin'),
('manager123', '$2a$08$6Y7QSNKu4oJ55nKNf4iBQueMiT7Hfg5.3p.UtNCB0EhX3chigXYau', 'manager'),
('viewer123', '$2a$08$6Y7QSNKu4oJ55nKNf4iBQueMiT7Hfg5.3p.UtNCB0EhX3chigXYau', 'viewer')
ON CONFLICT (username) DO NOTHING;

-- Права доступа
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO warehouse_control_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO warehouse_control_user;

COMMIT;