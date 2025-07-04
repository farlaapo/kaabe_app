-- ✅ Enable UUID Extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ✅ Create ENUM type for subscription status
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'subscription_status') THEN
        CREATE TYPE subscription_status AS ENUM ('active', 'expired', 'cancelled', 'pending');
    END IF;
END
$$;

-- ✅ Create Subscriptions Table
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    status subscription_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_subscription_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_subscription_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);


CREATE OR REPLACE PROCEDURE create_subscription(
    IN p_id UUID,
    IN p_user_id UUID,
    IN p_course_id UUID,
    IN p_started_at TIMESTAMPTZ,
    IN p_expires_at TIMESTAMPTZ,
    IN p_status subscription_status
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO subscriptions (
        id, user_id, course_id, started_at, expires_at, status, created_at, updated_at
    ) VALUES (
        p_id, p_user_id, p_course_id, p_started_at, p_expires_at, p_status, NOW(), NOW()
    );
END;
$$;

CREATE OR REPLACE PROCEDURE update_subscription(
    IN p_id UUID,
    IN p_started_at TIMESTAMPTZ,
    IN p_expires_at TIMESTAMPTZ,
    IN p_status subscription_status
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE subscriptions
    SET
        started_at = p_started_at,
        expires_at = p_expires_at,
        status = p_status,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_id AND deleted_at IS NULL;
END;
$$;


CREATE OR REPLACE PROCEDURE delete_subscription(
    IN p_id UUID
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE subscriptions
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = p_id;
END;
$$;


CREATE OR REPLACE FUNCTION get_subscription_by_id(p_id UUID)
RETURNS TABLE (
    id UUID,
    user_id UUID,
    course_id UUID,
    started_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    status subscription_status,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT
        subscriptions.id,
        subscriptions.user_id,
        subscriptions.course_id,
        subscriptions.started_at,
        subscriptions.expires_at,
        subscriptions.status,
        subscriptions.created_at,
        subscriptions.updated_at
    FROM subscriptions
    WHERE subscriptions.id = p_id AND subscriptions.deleted_at IS NULL;
END;
$$;



CREATE OR REPLACE FUNCTION get_all_subscriptions()
RETURNS TABLE (
    id UUID,
    user_id UUID,
    course_id UUID,
    started_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    status subscription_status,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT
        subscriptions.id,
        subscriptions.user_id,
        subscriptions.course_id,
        subscriptions.started_at,
        subscriptions.expires_at,
        subscriptions.status,
        subscriptions.created_at,
        subscriptions.updated_at
    FROM subscriptions
    WHERE subscriptions.deleted_at IS NULL
    ORDER BY subscriptions.created_at DESC;
END;
$$;

