CREATE TABLE IF NOT EXISTS flights (
    id TEXT PRIMARY KEY,
    flight_id TEXT NOT NULL,
    destination_from TEXT NOT NULL,
    destination_to TEXT NOT NULL
);