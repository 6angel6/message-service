CREATE TABLE Messages (
                          id SERIAL PRIMARY KEY,
                          content TEXT NOT NULL,
                          status VARCHAR(20) DEFAULT 'pending',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)