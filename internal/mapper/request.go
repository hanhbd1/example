package mapper

import (
	"example/internal/dto"
	"example/internal/models"
)

func ToModelFromRequest(m *dto.PersonnelCreateRequest, currentPersonnels ...*models.Personnel) (*models.Personnel, []*models.Education) {
	if m == nil {
		return nil, []*models.Education{}
	}
	var currentPersonnel *models.Personnel
	if len(currentPersonnels) > 0 {
		currentPersonnel = currentPersonnels[0]
	}
	p := FromPersonnelCreation(m.Person)
	if currentPersonnel != nil {
		p.Id = currentPersonnel.Id
		p.CreatedAt = currentPersonnel.CreatedAt
	}
	edus := FromEducationCreations(m.Educations, p.Id)
	return p, edus
}
