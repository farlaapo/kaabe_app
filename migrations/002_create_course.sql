-- Create extension for UUID if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create ENUM type for course status if not exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'course_status') THEN
        CREATE TYPE course_status AS ENUM ('draft', 'published', 'archived');
    END IF;
END
$$;

-- Create courses table if not exists
CREATE TABLE IF NOT EXISTS courses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    influencer_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price FLOAT,
    cover_image_url TEXT[],          -- array of text
    status course_status,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_courses_influencer FOREIGN KEY (influencer_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create the procedure with correct text[] parameter
CREATE OR REPLACE PROCEDURE create_course(
    IN p_id UUID,
    IN p_influencer_id UUID,
    IN p_title TEXT,
    IN p_description TEXT,
    IN p_price FLOAT8,
    IN p_cover_image_url TEXT[],
    IN p_status course_status
)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO courses (
        id, influencer_id, title, description, price, cover_image_url, status
    ) VALUES (
        p_id, p_influencer_id, p_title, p_description, p_price, p_cover_image_url, p_status
    );
END;
$$;



-- Create procedure: update_course
CREATE OR REPLACE PROCEDURE update_course(
    IN p_id UUID,
    IN p_title TEXT,
    IN p_description TEXT,
    IN p_price FLOAT8,
    IN p_cover_image_url TEXT[],
    IN p_status VARCHAR
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE courses
    SET
        title = p_title,
        description = p_description,
        price = p_price,
        cover_image_url = p_cover_image_url,
        status = p_status::course_status,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = p_id AND deleted_at IS NULL;
END;
$$;

-- Create procedure: soft delete course
CREATE OR REPLACE PROCEDURE delete_course(
    IN p_id UUID
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE courses
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = p_id;
END;
$$;

-- Create function: get all active courses
CREATE OR REPLACE FUNCTION get_all_courses()
RETURNS TABLE (
    id UUID,
    influencer_id UUID,
    title VARCHAR,
    description TEXT,
    price FLOAT,
    cover_image_url TEXT[],
    status VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT id, influencer_id, title, description, price, cover_image_url, status::text, created_at, updated_at
    FROM courses
    WHERE deleted_at IS NULL;
END;
$$;

-- Create function: get a single active course by ID
CREATE OR REPLACE FUNCTION get_course_by_id(p_course_id UUID)
RETURNS TABLE (
    id UUID,
    influencer_id UUID,
    title VARCHAR,
    description TEXT,
    price FLOAT,
    cover_image_url TEXT[],
    status VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE plpgsql AS $$
BEGIN
    RETURN QUERY
    SELECT 
        courses.id, 
        courses.influencer_id, 
        courses.title, 
        courses.description, 
        courses.price, 
        courses.cover_image_url, 
        courses.status::text, 
        courses.created_at, 
        courses.updated_at
    FROM courses
    WHERE courses.id = p_course_id AND courses.deleted_at IS NULL;
END;
$$;


