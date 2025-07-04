-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_ref VARCHAR(255) UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    subscription_id UUID NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    status VARCHAR(50) NOT NULL,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    -- Foreign key constraints
    CONSTRAINT fk_payment_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_payment_subscription FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
);


-- create payment
CREATE OR REPLACE PROCEDURE create_payment(
    IN p_id UUID,
    IN p_external_ref VARCHAR,
    IN p_user_id UUID,
    IN p_subscription_id UUID,
    IN p_amount DOUBLE PRECISION,
    IN p_status VARCHAR,
    IN p_processed_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO payments (id, external_ref, user_id, subscription_id, amount, status, processed_at, created_at, updated_at)
    VALUES (p_id, p_external_ref, p_user_id, p_subscription_id, p_amount, p_status, p_processed_at, NOW(), NOW());
END;
$$;

-- update payment 
CREATE OR REPLACE PROCEDURE update_payment(
    IN p_id UUID,
    IN p_external_ref VARCHAR,
    IN p_user_id UUID,
    IN p_subscription_id UUID,
    IN p_amount DOUBLE PRECISION,
    IN p_status VARCHAR,
    IN p_processed_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE payments
    SET external_ref = p_external_ref,
        user_id = p_user_id,
        subscription_id = p_subscription_id,
        amount = p_amount,
        status = p_status,
        processed_at = p_processed_at,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_id AND deleted_at IS NULL;
END;
$$;

-- delete payment 
CREATE OR REPLACE PROCEDURE delete_payment(
    IN p_id UUID
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE payments
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = p_id;
END;
$$;

-- create get payment by id
CREATE OR REPLACE FUNCTION get_payment_by_id(p_id UUID)
RETURNS TABLE (
    id UUID,
    external_ref VARCHAR,
    user_id UUID,
    subscription_id UUID,
    amount DOUBLE PRECISION,
    status VARCHAR,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT 
        payments.id,
        payments.external_ref,
        payments.user_id,
        payments.subscription_id,
        payments.amount,
        payments.status,
        payments.processed_at,
        payments.created_at,
        payments.updated_at
    FROM payments
    WHERE payments.id = p_id AND payments.deleted_at IS NULL;
END;
$$;


-- create get all payment by user id
CREATE OR REPLACE FUNCTION get_all_payments() 
RETURNS TABLE (
    id UUID,
    external_ref VARCHAR,
    user_id UUID,
    subscription_id UUID,
    amount DOUBLE PRECISION,
    status VARCHAR,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT
        payments.id,
        payments.external_ref,
        payments.user_id,
        payments.subscription_id,
        payments.amount,
        payments.status,
        payments.processed_at,
        payments.created_at,
        payments.updated_at
    FROM payments
    WHERE payments.deleted_at IS NULL
    ORDER BY payments.created_at DESC;
END;
$$;


-- Function to get payment by External Reference
CREATE OR REPLACE FUNCTION get_payment_by_external_ref(p_external_ref VARCHAR)
RETURNS TABLE (
    id UUID,
    external_ref VARCHAR,
    user_id UUID,
    subscription_id UUID,
    amount DOUBLE PRECISION,
    status VARCHAR,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT id, external_ref, user_id, subscription_id, amount, status, processed_at, created_at, updated_at
    FROM payments
    WHERE external_ref = p_external_ref AND deleted_at IS NULL;
END;
$$;
