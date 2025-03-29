CREATE TABLE IF NOT EXISTS flights (
    flight_id TEXT PRIMARY KEY,
    destination_from TEXT NOT NULL,
    destination_to TEXT NOT NULL
);