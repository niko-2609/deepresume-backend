CREATE EXTENSION IF NOT EXISTS vector; 

-- Create table for keyword vectors
CREATE TABLE IF NOT EXISTS keyword_vectors (
    id SERIAL PRIMARY KEY,
    source_type VARCHAR(10) NOT NULL, -- 'job' or 'resume'
    source_id VARCHAR(100) NOT NULL,
    text TEXT NOT NULL,
    keywords JSONB NOT NULL,
    vector_embedding vector(50), -- Dimension depends on your vector size
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(source_type, source_id)
);

-- Create index for faster vector similarity searches
CREATE INDEX IF NOT EXISTS keyword_vectors_source_type_idx ON keyword_vectors(source_type); 