CREATE TABLE grades (
    id uuid PRIMARY KEY,
    curriculum_id uuid NOT NULL REFERENCES curriculums(id),
    student_id uuid NOT NULL REFERENCES students(id),
    author_id uuid NOT NULL REFERENCES users(id),
    value int NOT NULL,
    comment varchar(255),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
