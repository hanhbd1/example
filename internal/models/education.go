package models

type Education struct {
	Id               string `gorm:"primaryKey;type:varchar(100)" json:"id"`
	CreatedAt        int64  `gorm:"type:bigint" json:"created_at"`
	PersonnelId      string `gorm:"type:varchar(100)" json:"personnel_id"`
	FromDate         string `gorm:"type:date" json:"from_date,omitempty"`
	ToDate           string `gorm:"type:date" json:"to_date,omitempty"`
	TrainingMethod   string `gorm:"type:varchar(255)" json:"training_method,omitempty"`
	EducationalLevel string `gorm:"type:varchar(255)" json:"educational_level,omitempty"`
	Major            string `gorm:"type:varchar(255)" json:"major,omitempty"`
	School           string `gorm:"type:varchar(255)" json:"school,omitempty"`
	//Personnel        *Personnel `gorm:"foreignKey:PersonnelId;references:Id"` // Many-to-One relationship (belongs to - use PersonnelID as foreign key) but I'm not recommended to use this
}

func (*Education) TableName() string {
	return "educations"
}
