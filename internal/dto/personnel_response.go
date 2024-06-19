package dto

type PersonnelResponse struct {
	Person     *Personnel   `json:"person"`
	Educations []*Education `json:"educations"`
}
