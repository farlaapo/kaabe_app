-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. Create ENUM for user roles (run once)
CREATE TYPE user_role AS ENUM ('admin', 'user', 'influencer');

-- 2. Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    wallet_id UUID,
    reset_token UUID,
    reset_token_expiry TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. Function: create_user
CREATE OR REPLACE FUNCTION public.create_user(
    p_email VARCHAR,
    p_password TEXT,
    p_first_name VARCHAR,
    p_last_name VARCHAR,
    p_role user_role DEFAULT 'user',
    p_wallet_id UUID DEFAULT NULL
)
RETURNS UUID
LANGUAGE plpgsql
AS $$
DECLARE
    new_user_id UUID;
BEGIN
    INSERT INTO users (id, email, password, first_name, last_name, role, wallet_id)
    VALUES (uuid_generate_v4(), p_email, p_password, p_first_name, p_last_name, p_role, p_wallet_id)
    RETURNING id INTO new_user_id;

    RETURN new_user_id;
END;
$$;

-- 4. Function: get_user_by_email
CREATE OR REPLACE FUNCTION get_user_by_email(p_email VARCHAR)
RETURNS TABLE (
    id UUID,
    email VARCHAR(255),
    password TEXT,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role user_role,
    wallet_id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE SQL
AS $$
    SELECT id, email, password, first_name, last_name, role, wallet_id, created_at, updated_at
    FROM users
    WHERE email = p_email;
$$;

-- 5. Function: get_user_by_id
CREATE OR REPLACE FUNCTION get_user_by_id(p_id UUID)
RETURNS TABLE (
    id UUID,
    email VARCHAR(255),
    password TEXT,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role user_role,
    wallet_id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT 
        id, email, password, first_name, last_name, role, wallet_id, created_at, updated_at
    FROM users
    WHERE id = p_id;
END;
$$;

-- 6. Function: get_all_users
CREATE OR REPLACE FUNCTION get_all_users()
RETURNS TABLE (
    id UUID,
    email VARCHAR(255),
    password TEXT,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role user_role,
    wallet_id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE SQL
AS $$
    SELECT id, email, password, first_name, last_name, role, wallet_id, created_at, updated_at
    FROM users
    ORDER BY created_at DESC;
$$;

-- 7. Procedure: update_user
CREATE OR REPLACE PROCEDURE update_user(
    IN p_id UUID,
    IN p_email VARCHAR,
    IN p_password TEXT,
    IN p_first_name VARCHAR,
    IN p_last_name VARCHAR,
    IN p_role user_role,
    IN p_wallet_id UUID,
    INOUT p_updated_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
DECLARE
    row_count INT;
BEGIN
    UPDATE users SET
        email = p_email,
        password = p_password,
        first_name = p_first_name,
        last_name = p_last_name,
        role = p_role,
        wallet_id = p_wallet_id,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_id;

    GET DIAGNOSTICS row_count = ROW_COUNT;

    IF row_count > 0 THEN
        SELECT updated_at INTO p_updated_at FROM users WHERE id = p_id;
    ELSE
        p_updated_at := NULL;
    END IF;
END;
$$;

-- 8. Function: delete_user
CREATE OR REPLACE FUNCTION delete_user(p_id UUID)
RETURNS INTEGER AS $$
DECLARE
    rows_deleted INTEGER;
BEGIN
    DELETE FROM users WHERE id = p_id;
    GET DIAGNOSTICS rows_deleted = ROW_COUNT;
    RETURN rows_deleted;
END;
$$ LANGUAGE plpgsql;

-- 9. Function: set_reset_token
CREATE OR REPLACE FUNCTION set_reset_token(
    p_email VARCHAR,
    p_token UUID,
    p_expiry TIMESTAMP
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE users
    SET reset_token = p_token,
        reset_token_expiry = p_expiry,
        updated_at = CURRENT_TIMESTAMP
    WHERE email = p_email;
END;
$$;

-- 10. Function: get_user_by_reset_token
CREATE OR REPLACE FUNCTION get_user_by_reset_token(p_token UUID)
RETURNS TABLE (
    id UUID,
    email VARCHAR(255),
    password TEXT,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role user_role,
    wallet_id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT
        users.id,
        users.email,
        users.password,
        users.first_name,
        users.last_name,
        users.role,
        users.wallet_id,
        users.created_at,
        users.updated_at
    FROM users
    WHERE users.reset_token = p_token
      AND users.reset_token_expiry > CURRENT_TIMESTAMP;
END;
$$;

-- 11. Function: update_user_password
CREATE OR REPLACE FUNCTION update_user_password(
    p_user_id UUID,
    p_password TEXT
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
    updated_count INTEGER;
BEGIN
    UPDATE users
    SET password = p_password,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_user_id;

    GET DIAGNOSTICS updated_count = ROW_COUNT;

    IF updated_count = 0 THEN
        RAISE EXCEPTION 'No user found with ID: %', p_user_id;
    END IF;
END;
$$;


-- 12. Function: clear_reset_token
CREATE OR REPLACE FUNCTION clear_reset_token(p_user_id UUID)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE users
    SET reset_token = NULL,
        reset_token_expiry = NULL,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_user_id;
END;
$$;

-- Procedures to create roles and permissions tables
CREATE OR REPLACE PROCEDURE create_roles_table()
LANGUAGE plpgsql AS $$
BEGIN
    CREATE TABLE IF NOT EXISTS roles (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) UNIQUE NOT NULL
    );
END;
$$;

CREATE OR REPLACE PROCEDURE create_permissions_table()
LANGUAGE plpgsql AS $$
BEGIN
    CREATE TABLE IF NOT EXISTS permissions (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) UNIQUE NOT NULL
    );
END;
$$;

CREATE OR REPLACE PROCEDURE create_user_roles_table()
LANGUAGE plpgsql AS $$
BEGIN
    CREATE TABLE IF NOT EXISTS user_roles (
        user_id UUID REFERENCES users(id) ON DELETE CASCADE,
        role_id INT REFERENCES roles(id) ON DELETE CASCADE,
        PRIMARY KEY (user_id, role_id)
    );
END;
$$;

CREATE OR REPLACE PROCEDURE create_user_permissions_table()
LANGUAGE plpgsql AS $$
BEGIN
    CREATE TABLE IF NOT EXISTS user_permissions (
        user_id UUID REFERENCES users(id) ON DELETE CASCADE,
        permission_id INT REFERENCES permissions(id) ON DELETE CASCADE,
        PRIMARY KEY (user_id, permission_id)
    );
END;
$$;
