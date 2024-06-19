package storage

import (
	"context"
	derrors "errors"

	"example/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (ds *DataStorage) CreateEducations(ctx context.Context, educations ...*models.Education) ([]*models.Education, error) {
	if err := ds.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).CreateInBatches(educations, 100).Error; err != nil {
		return nil, err
	}

	return educations, nil
}

func (ds *DataStorage) DeleteEducation(ctx context.Context, id string) error {
	Education := &models.Education{
		Id: id,
	}
	if err := ds.db.WithContext(ctx).Where(Education).Delete(Education).Error; err != nil {
		return err
	}

	return nil
}

func (ds *DataStorage) GetEducation(ctx context.Context, id string) (*models.Education, error) {
	Education := &models.Education{Id: id}
	err := ds.db.WithContext(ctx).Model(Education).Take(Education).Error
	if err != nil {
		if derrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return Education, nil
}

func (ds *DataStorage) ListEducations(ctx context.Context) ([]*models.Education, error) {
	res := []*models.Education{}
	if err := ds.db.WithContext(ctx).Model(&models.Education{}).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (ds *DataStorage) ListEducationsByPersonnelId(ctx context.Context, personnelId string) ([]*models.Education, error) {
	res := []*models.Education{}
	if err := ds.db.WithContext(ctx).Model(&models.Education{}).Where("personnel_id = ?", personnelId).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (ds *DataStorage) DeleteEducationsByPersonnelId(ctx context.Context, personnelId string) error {
	if err := ds.db.WithContext(ctx).Where("personnel_id = ?", personnelId).Delete(&models.Education{}).Error; err != nil {
		return err
	}
	return nil
}
