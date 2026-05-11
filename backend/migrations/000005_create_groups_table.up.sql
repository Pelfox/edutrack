CREATE TABLE groups (
    id uuid PRIMARY KEY,
    name varchar(255) NOT NULL,
    study_form study_form NOT NULL,
    admission_year int NOT NULL,
    specialty_id uuid NOT NULL REFERENCES specialties(id),
    created_at timestamptz NOT NULL
);
