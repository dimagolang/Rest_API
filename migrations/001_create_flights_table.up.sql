DROP TABLE IF EXISTS flights;

CREATE TABLE flights (
                         id SERIAL PRIMARY KEY,
                         destination_from VARCHAR(100) NOT NULL,
                         destination_to VARCHAR(100) NOT NULL,
                         delete_at BIGINT DEFAULT 0
);