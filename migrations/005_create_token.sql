-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create the tokens table
CREATE TABLE IF NOT EXISTS tokens (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Procedure: Create a token
CREATE OR REPLACE PROCEDURE create_token(
    IN p_id UUID,
    IN p_user_id UUID,
    IN p_token TEXT,
    IN p_created_at TIMESTAMP
    IN p_expires_at TIMESTAMP
    IN p_updated_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO tokens (
        id, user_id, token, expires_at, created_at, updated_at
    ) VALUES (
        p_id, p_user_id, p_token, p_expires_at, NOW(), NOW()
    );
END;
$$;


CREATE OR REPLACE FUNCTION get_token_by_token(p_token TEXT)
RETURNS TABLE (
    id UUID,
    user_id UUID,
    token TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        id, user_id, token, expires_at, created_at, updated_at, deleted_at
    FROM 
        tokens
    WHERE 
        token = p_token AND deleted_at IS NULL;
END;
$$;



