CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY,
    parking_id UUID NOT NULL,
    active BOOLEAN NOT NULL,
    valid_from DATE NOT NULL,
    valid_to DATE,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    week_days SMALLINT[] NOT NULL,
    cost BIGINT
);

ALTER TABLE schedules
    ADD CONSTRAINT fk_schedules_parking
    FOREIGN KEY (parking_id)
    REFERENCES parkings(id)
    ON DELETE RESTRICT;