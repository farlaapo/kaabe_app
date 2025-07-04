-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the ratings table
CREATE TABLE IF NOT EXISTS ratings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    score INT NOT NULL CHECK (score >= 1 AND score <= 5),
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_ratings_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_ratings_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Procedure: Create a new rating with explicit id
CREATE OR REPLACE PROCEDURE create_rating(
    IN p_id UUID,
    IN p_user_id UUID,
    IN p_course_id UUID,
    IN p_score INT,
    IN p_comment TEXT
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO ratings (id, user_id, course_id, score, comment, created_at, updated_at)
    VALUES (p_id, p_user_id, p_course_id, p_score, p_comment, NOW(), NOW());
END;
$$;

-- Procedure: Update an existing rating
CREATE OR REPLACE PROCEDURE update_rating(
    IN p_id UUID,
    IN p_score INT,
    IN p_comment TEXT
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE ratings
    SET score = p_score,
        comment = p_comment,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_id AND deleted_at IS NULL;
END;
$$;

-- Procedure: Soft delete a rating
CREATE OR REPLACE PROCEDURE delete_rating(
    IN p_id UUID
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE ratings
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = p_id;
END;
$$;

-- Function: Get a single active rating by ID
CREATE OR REPLACE FUNCTION get_rating_by_id(p_id UUID)
RETURNS TABLE (
    id UUID,
    user_id UUID,
    course_id UUID,
    score INT,
    comment TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT 
        ratings.id,
        ratings.user_id,
        ratings.course_id,
        ratings.score,
        ratings.comment,
        ratings.created_at,
        ratings.updated_at
    FROM ratings
    WHERE ratings.id = p_id AND ratings.deleted_at IS NULL;
END;
$$;


-- Function: Get all active ratings
CREATE OR REPLACE FUNCTION get_all_ratings()
RETURNS TABLE (
    id UUID,
    user_id UUID,
    course_id UUID,
    score INT,
    comment TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT
        ratings.id,
        ratings.user_id,
        ratings.course_id,
        ratings.score,
        ratings.comment,
        ratings.created_at,
        ratings.updated_at
    FROM ratings
    WHERE ratings.deleted_at IS NULL
    ORDER BY ratings.created_at DESC;
END;
$$;


