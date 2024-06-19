package dto

type EducationCreation struct {
	FromDate         string `json:"from_date,omitempty"`
	ToDate           string `json:"to_date,omitempty"`
	TrainingMethod   string `json:"training_method,omitempty"`
	EducationalLevel string `json:"educational_level,omitempty"`
	Major            string `json:"major,omitempty"`
	School           string `json:"school,omitempty"`
}
