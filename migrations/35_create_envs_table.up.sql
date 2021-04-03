CREATE TABLE IF NOT EXISTS env_values(
	id INTEGER PRIMARY KEY,
	application_id INTEGER NOT NULL,
        key TEXT NOT NULL,
        value TEXT NOT NULL
);

