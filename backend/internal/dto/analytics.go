package dto

// AnalyticsGradeDistribution содержит количество оценок по одному значению.
type AnalyticsGradeDistribution struct {
	Value int `json:"value"`
	Count int `json:"count"`
}

// AnalyticsSubjectAverage содержит средний балл по одной дисциплине.
type AnalyticsSubjectAverage struct {
	SubjectID    string  `json:"subject_id"`
	SubjectTitle string  `json:"subject_title"`
	AverageGrade float64 `json:"average_grade"`
	GradesCount  int     `json:"grades_count"`
}

// AnalyticsOverview содержит сводные показатели админ-панели.
type AnalyticsOverview struct {
	StudentsCount     int                          `json:"students_count"`
	TeachersCount     int                          `json:"teachers_count"`
	GroupsCount       int                          `json:"groups_count"`
	SpecialtiesCount  int                          `json:"specialties_count"`
	SubjectsCount     int                          `json:"subjects_count"`
	CurriculumsCount  int                          `json:"curriculums_count"`
	GradesCount       int                          `json:"grades_count"`
	AverageGrade      *float64                     `json:"average_grade"`
	GradeDistribution []AnalyticsGradeDistribution `json:"grade_distribution"`
	SubjectAverages   []AnalyticsSubjectAverage    `json:"subject_averages"`
}
