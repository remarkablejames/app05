-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('superuser', 'admin', 'instructor', 'student');

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       email VARCHAR(255) NOT NULL UNIQUE,
                       hashed_password VARCHAR(255) NOT NULL,
                       first_name VARCHAR(100) NOT NULL,
                       last_name VARCHAR(100) NOT NULL,
                       profile_picture_url VARCHAR(255),
                       role user_role NOT NULL,
                       active BOOLEAN NOT NULL DEFAULT true,
                       email_verified BOOLEAN NOT NULL DEFAULT true,
                       subscribed_to_newsletter BOOLEAN NOT NULL DEFAULT false,
                       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                       last_login_at TIMESTAMP WITH TIME ZONE,
                       password_reset_token VARCHAR(255),
                       reset_token_expires_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for better query performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_active ON users(active);

-- Trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Seed a special "Deleted User" immediately after table creation
INSERT INTO users (
    id,
    email,
    hashed_password,
    first_name,
    last_name,
    profile_picture_url,
    role,
    active,
    email_verified,
    created_at,
    updated_at,
    last_login_at,
    password_reset_token,
    reset_token_expires_at
) VALUES (
             uuid_generate_v4(),
             'deleted_user@somolabs.com',
             'fd9a89976669a5a8821a108c151a99ffc5823bac',
             'Deleted',
             'User',
             'This is a placeholder account for deleted users.',
             'student',
             false,
             false,
             CURRENT_TIMESTAMP,
             CURRENT_TIMESTAMP,
             '2023-01-01 00:00:00',
             NULL,
             NULL
         );