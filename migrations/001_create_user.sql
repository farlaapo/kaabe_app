-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create ENUM type for user_role if not exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('user', 'admin', 'influencer');
    END IF;
END;
$$ LANGUAGE plpgsql;


-- 1. Create users table (run once)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    wallet_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE OR REPLACE FUNCTION public.create_user(
    p_email character varying,
    p_password text,
    p_first_name character varying,
    p_last_name character varying,
    p_role user_role DEFAULT 'user',
    p_wallet_id uuid DEFAULT NULL::uuid
)
RETURNS uuid
LANGUAGE plpgsql
AS $$
DECLARE
    new_user_id uuid;
BEGIN
    INSERT INTO users (id, email, password, first_name, last_name, role, wallet_id)
    VALUES (uuid_generate_v4(), p_email, p_password, p_first_name, p_last_name, p_role, p_wallet_id)
    RETURNING id INTO new_user_id;

    RETURN new_user_id;
END;
$$;



-- 2. Function: get_user_by_email
CREATE OR REPLACE FUNCTION get_user_by_email(p_email VARCHAR)
RETURNS TABLE (
    id UUID,
    email VARCHAR,
    password TEXT,
    first_name VARCHAR,
    last_name VARCHAR,
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


-- 3. Function: get_user_by_id
CREATE OR REPLACE FUNCTION get_user_by_id(p_id UUID)
RETURNS TABLE (
    id UUID,
    email TEXT,
    password TEXT,
    first_name TEXT,
    last_name TEXT,
    role TEXT,
    wallet_id UUID,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT 
        id,
        email,
        password,
        first_name,
        last_name,
        role,
        wallet_id,
        created_at,
        updated_at
    FROM users
    WHERE id = p_id AND deleted_at IS NULL;
END;
$$;



-- 4. Function: get_all_users
CREATE OR REPLACE FUNCTION get_all_users()
RETURNS TABLE (
    id UUID,
    email TEXT,
    password TEXT,
    first_name TEXT,
    last_name TEXT,
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


-- 5. Procedure: update_user
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


-- 6. FUNCTION: delete_user
CREATE OR REPLACE FUNCTION delete_user(p_id UUID) RETURNS INTEGER AS $$
DECLARE
    rows_deleted INTEGER;
BEGIN
    DELETE FROM users WHERE id = p_id;
    GET DIAGNOSTICS rows_deleted = ROW_COUNT;
    RETURN rows_deleted;
END;
$$ LANGUAGE plpgsql;



--- Create create_roles table
CREATE OR REPLACE PROCEDURE create_roles_table()
LANGUAGE plpgsql AS $$
BEGIN
    CREATE TABLE IF NOT EXISTS roles (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) UNIQUE NOT NULL
    );
END;
$$;


-- Create create_permissions table
CREATE OR REPLACE PROCEDURE create_permissions_table()
LANGUAGE plpgsql AS $$
BEGIN
    CREATE TABLE IF NOT EXISTS permissions (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) UNIQUE NOT NULL
    );
END;
$$;


-- Create a table to store user roles
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

----- Create create_user_permissions table
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
