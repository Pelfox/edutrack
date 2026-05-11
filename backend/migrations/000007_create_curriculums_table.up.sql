CREATE TABLE curriculums (
    id uuid PRIMARY KEY,
    hours int NOT NULL,
    semester int NOT NULL,
    report_type curriculum_report_type NOT NULL,
    subject_id uuid NOT NULL REFERENCES subjects(id),
    group_id uuid NOT NULL REFERENCES groups(id),
    lead_by uuid NOT NULL REFERENCES users(id)
);
