BEGIN;

-- Drop triggers first
DROP TRIGGER IF EXISTS update_balances_uploaded_at ON balance;
DROP TRIGGER IF EXISTS update_orders_uploaded_at ON orders;

-- Drop function
DROP FUNCTION IF EXISTS update_uploaded_at_column();
DROP function IF EXISTS withdraw(p_user_id INT, p_order_number VARCHAR, p_sum numeric(10, 2));

-- Drop indexes
DROP INDEX IF EXISTS idx_balances_user_id;
DROP INDEX IF EXISTS idx_withdrawals_processed_at;
DROP INDEX IF EXISTS idx_withdrawals_user_id;
DROP INDEX IF EXISTS idx_orders_uploaded_at;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_orders_order_number;
DROP INDEX IF EXISTS idx_users_login;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS withdrawals;
DROP TABLE IF EXISTS balance;
DROP TABLE IF EXISTS orders;

-- Drop custom type
DROP TYPE IF EXISTS order_status_type;

-- Drop users table last
DROP TABLE IF EXISTS users;

COMMIT;
