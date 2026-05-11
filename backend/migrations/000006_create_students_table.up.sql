CREATE TABLE students (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id),
    group_id uuid NOT NULL REFERENCES groups(id),
    last_name varchar(255) NOT NULL,
    first_name varchar(255) NOT NULL,
    middle_name varchar(255),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
