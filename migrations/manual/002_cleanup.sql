BEGIN;

REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM warehouse_control_user;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM warehouse_control_user;

DROP INDEX IF EXISTS idx_items_history_item_id;
DROP INDEX IF EXISTS idx_items_history_user_id;
DROP INDEX IF EXISTS idx_items_history_changed_at;

DROP TABLE IF EXISTS items_history;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS users;

COMMIT;