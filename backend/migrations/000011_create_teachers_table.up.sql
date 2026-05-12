CREATE TABLE teachers (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    last_name varchar(255) NOT NULL,
    first_name varchar(255) NOT NULL,
    middle_name varchar(255),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
