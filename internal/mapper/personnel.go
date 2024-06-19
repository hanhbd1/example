package mapper

import (
	"time"

	"example/internal/dto"
	"example/internal/models"

	"github.com/google/uuid"
)

func ToPersonnelDTO(m *models.Personnel) *dto.Personnel {
	if m == nil {
		return nil
	}
	recordDTO := &dto.Personnel{
		Id:            m.Id,
		CreatedAt:     m.CreatedAt,
		FullName:      m.FullName,
		Dob:           m.Dob,
		Gender:        m.Gender,
		PlaceOfOrigin: m.PlaceOfOrigin,
		PlaceOfBirth:  m.PlaceOfBirth,
		Ethnicity:     m.Ethnicity,
		Religion:      m.Religion,
		Nationality:   m.Nationality,
		MaritalStatus: m.MaritalStatus,
		TaxCode:       m.TaxCode,
	}
	return recordDTO
}

func FromPersonnelCreation(m dto.PersonnelCreation) *models.Personnel {
	recordModel := &models.Personnel{
		Id:            uuid.Must(uuid.NewRandom()).String(),
		CreatedAt:     time.Now().Unix(),
		FullName:      m.FullName,
		Dob:           m.Dob,
		Gender:        m.Gender,
		PlaceOfOrigin: m.PlaceOfOrigin,
		PlaceOfBirth:  m.PlaceOfBirth,
		Ethnicity:     m.Ethnicity,
		Religion:      m.Religion,
		Nationality:   m.Nationality,
		MaritalStatus: m.MaritalStatus,
		TaxCode:       m.TaxCode,
	}
	return recordModel
}

func ToPersonnelModel(m *dto.Personnel) *models.Personnel {
	recordModel := &models.Personnel{
		Id:            m.Id,
		CreatedAt:     m.CreatedAt,
		FullName:      m.FullName,
		Dob:           m.Dob,
		Gender:        m.Gender,
		PlaceOfOrigin: m.PlaceOfOrigin,
		PlaceOfBirth:  m.PlaceOfBirth,
		Ethnicity:     m.Ethnicity,
		Religion:      m.Religion,
		Nationality:   m.Nationality,
		MaritalStatus: m.MaritalStatus,
		TaxCode:       m.TaxCode,
	}
	return recordModel
}
