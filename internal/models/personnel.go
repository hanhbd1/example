package models

type Personnel struct {
	Id            string `gorm:"primaryKey;type:varchar(100)" json:"id"`
	CreatedAt     int64  `gorm:"type:bigint" json:"created_at"`
	FullName      string `gorm:"type:varchar(255);not null" json:"full_name"`
	Dob           string `gorm:"type:date;not null" json:"dob"`
	Gender        string `gorm:"type:varchar(50);not null" json:"gender"`
	PlaceOfOrigin string `gorm:"type:varchar(255);not null" json:"place_of_origin"`
	PlaceOfBirth  string `gorm:"type:varchar(255);not null" json:"place_of_birth"`
	Ethnicity     string `gorm:"type:varchar(255)" json:"ethnicity"`
	Religion      string `gorm:"type:varchar(255)" json:"religion"`
	Nationality   string `gorm:"type:varchar(255)" json:"nationality"`
	MaritalStatus string `gorm:"type:varchar(255);not null" json:"marital_status"`
	TaxCode       string `gorm:"type:varchar(255)" json:"tax_code"`
	//Educations    []*Education `gorm:"foreignKey:PersonnelId"` // One-to-Many relationship (has many - use PersonnelID as foreign key) but I'm not recommended to use this
}

func (*Personnel) TableName() string {
	return "personnels"
}
