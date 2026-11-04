CREATE TABLE IF NOT EXISTS parking_sessions (
    id UUID PRIMARY KEY,
    parking_id UUID NOT NULL,
    state_number VARCHAR(9) NOT NULL,
    entered_at TIMESTAMPTZ NOT NULL,
    exited_at TIMESTAMPTZ,
    cost BIGINT,
    paid BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_parking_sessions_parking_id
    ON parking_sessions (parking_id);

CREATE INDEX idx_parking_sessions_state_number
    ON parking_sessions (state_number);

CREATE INDEX idx_parking_sessions_active
    ON parking_sessions (parking_id, state_number)
    WHERE exited_at IS NULL;