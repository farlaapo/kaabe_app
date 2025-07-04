
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the lessons table
CREATE TABLE IF NOT EXISTS lessons (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    course_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    video_url TEXT[],  -- Store multiple video URLs
    lesson_order INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_lessons_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
);

-- Create new procedure with `id` as a parameter
CREATE OR REPLACE PROCEDURE create_lesson(
    IN p_id UUID,
    IN p_course_id UUID,
    IN p_title VARCHAR(255),
    IN p_video_url TEXT[],
    IN p_order INT
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO lessons (id, course_id, title, video_url, lesson_order, created_at, updated_at)
    VALUES (p_id, p_course_id, p_title, p_video_url, p_order, NOW(), NOW());
END;
$$;


-- Procedure: Update an existing lesson
CREATE OR REPLACE PROCEDURE update_lesson(
    IN p_id UUID,
    IN p_title VARCHAR(255),
    IN p_video_url TEXT[],
    IN p_order INT
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE lessons
    SET
        title = p_title,
        video_url = p_video_url,
        lesson_order = p_order,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_id AND deleted_at IS NULL;
END;
$$;

-- Procedure: Soft delete a lesson
CREATE OR REPLACE PROCEDURE delete_lesson(
    IN p_id UUID
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE lessons
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = p_id;
END;
$$;


-- Function: Get a single active lesson by lesson ID
CREATE OR REPLACE FUNCTION get_lesson_by_id(p_lesson_id UUID)
RETURNS TABLE (
    id UUID,
    course_id UUID,
    title TEXT,
    video_url TEXT[],
    lesson_order INT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT 
        lessons.id,
        lessons.course_id,
        lessons.title::TEXT,
        lessons.video_url,
        lessons.lesson_order,
        lessons.created_at,
        lessons.updated_at
    FROM lessons
    WHERE lessons.id = p_lesson_id AND lessons.deleted_at IS NULL;
END;
$$;


-- Function: Get all active lessons
CREATE OR REPLACE FUNCTION get_all_lessons()
RETURNS TABLE (
    id UUID,
    course_id UUID,
    title TEXT,
    video_url TEXT[],
    lesson_order INT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT
        lessons.id,
        lessons.course_id,
        lessons.title,
        lessons.video_url,
        lessons.lesson_order,
        lessons.created_at,
        lessons.updated_at
    FROM lessons
    WHERE lessons.deleted_at IS NULL
    ORDER BY lessons.created_at DESC;
END;
$$;









