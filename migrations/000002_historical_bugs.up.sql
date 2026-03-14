-- Enable pgvector extension and create historical_bugs table

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS historical_bugs (
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    priority TEXT,
    category TEXT,
    embedding vector(3072)
);
