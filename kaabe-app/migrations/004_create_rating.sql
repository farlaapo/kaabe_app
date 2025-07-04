-- Create extension for UUID
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

    CONSTRAINT fk_rating_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_rating_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Procedure to add a new rating
CREATE OR REPLACE PROCEDURE create_rating(
    p_user_id UUID,
    p_course_id UUID,
    p_score INT,
    p_comment TEXT
)
LANGUAGE plpgsql AS $$
BEGIN 
    INSERT INTO ratings(user_id, course_id, score, comment)
    VALUES (p_user_id, p_course_id, p_score, p_comment);
END;
$$;

-- Procedure to update a rating
CREATE OR REPLACE PROCEDURE update_rating(
    p_id UUID,
    p_score INT,
    p_comment TEXT
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

-- Procedure to soft delete a rating
CREATE OR REPLACE PROCEDURE delete_rating(p_id UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE ratings
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = p_id;
END;
$$;

-- Function to get a rating by its ID
CREATE OR REPLACE FUNCTION get_rating_by_id(p_id UUID)
RETURNS TABLE (
    id UUID,
    user_id UUID,
    course_id UUID,
    score INT,
    comment TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
) AS $$
BEGIN 
    RETURN QUERY
    SELECT id, user_id, course_id, score, comment, created_at, updated_at
    FROM ratings
    WHERE id = p_id AND deleted_at IS NULL;
END;
$$ LANGUAGE plpgsql;

-- Function to get all ratings
CREATE OR REPLACE FUNCTION get_all_ratings()
RETURNS TABLE (
    id UUID,
    user_id UUID,
    course_id UUID,
    score INT,
    comment TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    SELECT id, user_id, course_id, score, comment, created_at, updated_at
    FROM ratings
    WHERE deleted_at IS NULL;
END;
$$ LANGUAGE plpgsql;