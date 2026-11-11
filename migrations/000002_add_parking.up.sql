CREATE TABLE IF NOT EXISTS parkings (
    id UUID PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE,
    address TEXT NOT NULL,
    capacity INTEGER NOT NULL CHECK (capacity > 0)
);

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_parkings_address_trgm
    ON parkings USING GIN (address gin_trgm_ops);

ALTER TABLE parking_sessions
    ADD CONSTRAINT fk_parking_sessions_parking
    FOREIGN KEY (parking_id)
    REFERENCES parkings(id)
    ON DELETE RESTRICT;