CREATE TABLE subjects (
    id uuid PRIMARY KEY,
    title varchar(255) NOT NULL,
    created_at timestamptz NOT NULL
);
