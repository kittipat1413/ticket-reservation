-- 202504221346_create_helper_function.down.sql

DROP EXTENSION IF EXISTS "pgcrypto";

DROP FUNCTION IF EXISTS update_updated_at_column;