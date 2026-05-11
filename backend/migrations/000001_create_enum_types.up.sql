CREATE TYPE user_role AS ENUM (
    'administrator',
    'teacher',
    'student'
);

CREATE TYPE study_form AS ENUM (
    'full_time',
    'evening',
    'extramural'
);

CREATE TYPE curriculum_report_type AS ENUM (
    'exam',
    'test',
    'diff_test'
);
