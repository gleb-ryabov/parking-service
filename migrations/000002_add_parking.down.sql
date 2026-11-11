ALTER TABLE parking_sessions
DROP CONSTRAINT fk_parking_sessions_parking;

DROP TABLE IF EXISTS parkings;