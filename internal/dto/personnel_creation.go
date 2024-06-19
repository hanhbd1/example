package dto

type PersonnelCreation struct {
	FullName      string `json:"full_name"`
	Dob           string `json:"dob"`
	Gender        string `json:"gender"`
	PlaceOfOrigin string `json:"place_of_origin"`
	PlaceOfBirth  string `json:"place_of_birth"`
	Ethnicity     string `json:"ethnicity"`
	Religion      string `json:"religion"`
	Nationality   string `json:"nationality"`
	MaritalStatus string `json:"marital_status"`
	TaxCode       string `json:"tax_code"`
}
