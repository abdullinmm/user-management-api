-- Drop trigger and function
DROP TRIGGER IF EXISTS trigger_create_user_balance ON users;
DROP FUNCTION IF EXISTS create_user_balance();

-- Drop tables in reverse order (respect foreign keys)
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS balances;
DROP TABLE IF EXISTS user_tasks;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;

-- Drop extension if needed
-- DROP EXTENSION IF EXISTS "uuid-ossp";
