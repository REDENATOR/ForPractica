-- Create users table (GORM model)
CREATE TABLE IF NOT EXISTS users (
    -- gorm.Model fields
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    fio TEXT NOT NULL,
    group_of_students TEXT,
    phone_number TEXT
);

-- Indexes for common queries
CREATE INDEX idx_users_group_of_students ON users(group_of_students);
CREATE INDEX idx_users_fio ON users(fio);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Function to update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();