package mapper

import (
	"time"

	"example/internal/dto"
	"example/internal/models"

	"github.com/google/uuid"
)

func ToEducationDTO(m *models.Education) *dto.Education {
	if m == nil {
		return nil
	}
	recordDTO := &dto.Education{
		Id:               m.Id,
		CreatedAt:        m.CreatedAt,
		PersonnelId:      m.PersonnelId,
		FromDate:         m.FromDate,
		ToDate:           m.ToDate,
		TrainingMethod:   m.TrainingMethod,
		EducationalLevel: m.EducationalLevel,
		Major:            m.Major,
		School:           m.School,
	}
	return recordDTO
}

func ToEducationDTOs(m []*models.Education) []*dto.Education {
	recordDTOs := make([]*dto.Education, 0)
	for _, v := range m {
		recordDTOs = append(recordDTOs, ToEducationDTO(v))
	}
	return recordDTOs
}

func FromEducationCreation(m dto.EducationCreation, personnelId string) *models.Education {
	recordModel := &models.Education{
		Id:               uuid.Must(uuid.NewRandom()).String(),
		CreatedAt:        time.Now().Unix(),
		PersonnelId:      personnelId,
		FromDate:         m.FromDate,
		ToDate:           m.ToDate,
		TrainingMethod:   m.TrainingMethod,
		EducationalLevel: m.EducationalLevel,
		Major:            m.Major,
		School:           m.School,
	}
	return recordModel
}

func FromEducationCreations(m []dto.EducationCreation, personnelId string) []*models.Education {
	recordModels := make([]*models.Education, 0)
	for _, v := range m {
		recordModels = append(recordModels, FromEducationCreation(v, personnelId))

	}
	return recordModels
}

func ToEducationModel(m *dto.Education) *models.Education {
	recordModel := &models.Education{
		Id:               m.Id,
		CreatedAt:        m.CreatedAt,
		PersonnelId:      m.PersonnelId,
		FromDate:         m.FromDate,
		ToDate:           m.ToDate,
		TrainingMethod:   m.TrainingMethod,
		EducationalLevel: m.EducationalLevel,
		Major:            m.Major,
		School:           m.School,
	}
	return recordModel
}
