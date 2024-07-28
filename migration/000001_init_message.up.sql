CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Messages (
                          id UUID DEFAULT uuid_generate_v4(),
                          content TEXT NOT NULL,
                          status VARCHAR(20) DEFAULT 'pending',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          PRIMARY KEY(id)
)