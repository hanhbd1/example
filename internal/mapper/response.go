package mapper

import (
	"example/internal/dto"
	"example/internal/models"
)

func ToDTOFromModels(p *models.Personnel, educations []*models.Education) *dto.PersonnelResponse {
	personDTO := &dto.PersonnelResponse{
		Person:     ToPersonnelDTO(p),
		Educations: ToEducationDTOs(educations),
	}
	return personDTO
}
