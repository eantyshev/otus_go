CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE appointment (
    uuid uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    summary text NOT NULL,
    description text DEFAULT '',
    time_start timestamp NOT NULL,
    time_end timestamp NOT NULL,
    owner text NOT NULL
);

CREATE INDEX start_idx ON appointment (owner, time_start);
