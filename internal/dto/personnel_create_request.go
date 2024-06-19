package dto

type PersonnelCreateRequest struct {
	Person     PersonnelCreation   `json:"person"`
	Educations []EducationCreation `json:"educations"`
}
