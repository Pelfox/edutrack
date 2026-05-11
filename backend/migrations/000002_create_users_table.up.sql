CREATE TABLE users (
    id uuid PRIMARY KEY,
    email varchar(255) NOT NULL,
    password_hash text NOT NULL,
    role user_role NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
