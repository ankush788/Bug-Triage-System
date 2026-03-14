-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create bugs table
CREATE TABLE IF NOT EXISTS bugs (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    reporter_id BIGINT NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN',
    priority VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN',
    category VARCHAR(100) NOT NULL DEFAULT 'UNCLASSIFIED',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- Create index on reporter_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_bugs_reporter_id ON bugs(reporter_id);

-- Create index on status for filtering
CREATE INDEX IF NOT EXISTS idx_bugs_status ON bugs(status);

-- Create index on priority for filtering
CREATE INDEX IF NOT EXISTS idx_bugs_priority ON bugs(priority);


