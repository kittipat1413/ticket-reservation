-- 202504221346_create_helper_function.up.sql

-- Enable UUID generation extension (one of these, depending on your setup)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create a helper function to update the updated_at column
CREATE FUNCTION update_updated_at_column() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
  BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
  END;
$$;