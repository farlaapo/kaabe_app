-- 1️⃣ Enable UUID Extension (for uuid_generate_v4)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2️⃣ ENUM: Withdrawal Status
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'withdrawal_status') THEN
        CREATE TYPE withdrawal_status AS ENUM ('pending', 'approved', 'rejected');
    END IF;
END
$$;

-- 3️⃣ TABLE: Withdrawals
CREATE TABLE IF NOT EXISTS withdrawals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    influencer_id UUID NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    status withdrawal_status NOT NULL DEFAULT 'pending',
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_withdrawals_influencer FOREIGN KEY (influencer_id) REFERENCES users(id) ON DELETE CASCADE
);



CREATE OR REPLACE PROCEDURE create_withdrawal(
    IN p_id UUID,
    IN p_influencer_id UUID,
    IN p_amount NUMERIC,
    IN p_status withdrawal_status,
    IN p_requested_at TIMESTAMPTZ,
    IN p_processed_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO withdrawals (
        id, influencer_id, amount, status, requested_at, processed_at, created_at, updated_at
    ) VALUES (
        p_id, p_influencer_id, p_amount, p_status, p_requested_at, p_processed_at, NOW(), NOW()
    );
END;
$$;

CREATE OR REPLACE PROCEDURE update_withdrawal(
    IN p_id UUID,
    IN p_amount NUMERIC,
    IN p_status withdrawal_status,
    IN p_processed_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE withdrawals
    SET
        amount = p_amount,
        status = p_status,
        processed_at = p_processed_at,
        updated_at = NOW()
    WHERE id = p_id AND deleted_at IS NULL;
END;
$$;

CREATE OR REPLACE PROCEDURE delete_withdrawal(
    IN p_id UUID
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE withdrawals
    SET deleted_at = NOW()
    WHERE id = p_id;
END;
$$;


CREATE OR REPLACE FUNCTION get_withdrawal_by_id(p_id UUID)
RETURNS TABLE (
    id UUID,
    influencer_id UUID,
    amount NUMERIC,
    status withdrawal_status,
    requested_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT
        withdrawals.id,
        withdrawals.influencer_id,
        withdrawals.amount,
        withdrawals.status,
        withdrawals.requested_at,
        withdrawals.processed_at,
        withdrawals.created_at,
        withdrawals.updated_at
    FROM withdrawals
    WHERE withdrawals.id = p_id AND withdrawals.deleted_at IS NULL;
END;
$$;


CREATE OR REPLACE FUNCTION get_all_withdrawals()
RETURNS TABLE (
    id UUID,
    influencer_id UUID,
    amount NUMERIC,
    status withdrawal_status,
    requested_at TIMESTAMPTZ,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT
        w.id,
        w.influencer_id,
        w.amount,
        w.status,
        w.requested_at,
        w.processed_at,
        w.created_at,
        w.updated_at
    FROM withdrawals w
    WHERE w.deleted_at IS NULL
    ORDER BY w.created_at DESC;
END;
$$;


