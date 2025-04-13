-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create posts table
CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       user_id UUID NOT NULL REFERENCES users(id),
                       title VARCHAR(255) NOT NULL,
                       content TEXT NOT NULL,
                       excerpt VARCHAR(500),
                       status VARCHAR(20) NOT NULL DEFAULT 'draft', -- draft, published, archived
                       slug VARCHAR(255) UNIQUE,
                       view_count INTEGER DEFAULT 0,
                       published_at TIMESTAMP WITH TIME ZONE,
                       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better query performance
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_published_at ON posts(published_at);
CREATE INDEX idx_posts_slug ON posts(slug);

-- Trigger to automatically update updated_at timestamp
CREATE TRIGGER update_posts_updated_at
    BEFORE UPDATE ON posts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Seed some dummy data
-- First, get a user ID to use in our seed data
DO $$
DECLARE
random_user_id UUID;
BEGIN
    -- Get a random user ID (excluding the deleted user)
SELECT id INTO random_user_id FROM users
WHERE email != 'deleted_user@somolabs.com' AND active = true
    LIMIT 1;

-- If no user found, use a new UUID
IF random_user_id IS NULL THEN
        random_user_id := uuid_generate_v4();

        -- Insert a new user if needed
INSERT INTO users (
    id,
    email,
    hashed_password,
    first_name,
    last_name,
    role,
    active,
    email_verified
) VALUES (
             random_user_id,
             'test_user@example.com',
             'hashed_password_placeholder',
             'Test',
             'User',
             'student',
             true,
             true
         );
END IF;

    -- Insert dummy posts
INSERT INTO posts (user_id, title, content, excerpt, status, slug, published_at) VALUES
                                                                                     (random_user_id, 'Getting Started with Go', 'This is a detailed guide about how to get started with Go programming language.', 'A beginners guide to Go', 'published', 'getting-started-with-go', CURRENT_TIMESTAMP - INTERVAL '5 days'),

                                                                                     (random_user_id, 'Advanced PostgreSQL Tips', 'PostgreSQL offers many advanced features that developers often overlook.', 'Discover advanced PostgreSQL features', 'published', 'advanced-postgresql-tips', CURRENT_TIMESTAMP - INTERVAL '3 days'),

                                                                                     (random_user_id, 'Building RESTful APIs with Go', 'Learn how to build robust RESTful APIs using Go and popular frameworks.', 'Guide to building APIs with Go', 'published', 'building-restful-apis-with-go', CURRENT_TIMESTAMP - INTERVAL '1 day'),

                                                                                     (random_user_id, 'Upcoming Go Features', 'Exploring the upcoming features in Go 2.0 and what they mean for developers.', 'The future of Go programming', 'draft', 'upcoming-go-features', NULL),

                                                                                     (random_user_id, 'Database Migration Strategies', 'Effective strategies for handling database migrations in production environments.', 'Manage database migrations effectively', 'published', 'database-migration-strategies', CURRENT_TIMESTAMP);
END $$;