package storage

import (
	"context"
	derrors "errors"

	"example/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (ds *DataStorage) CreatePersonnel(ctx context.Context, personnel *models.Personnel) (*models.Personnel, error) {
	if err := ds.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(personnel).Error; err != nil {
		return nil, err
	}

	return personnel, nil
}

func (ds *DataStorage) UpdatePersonnel(ctx context.Context, personnel *models.Personnel) (*models.Personnel, error) {
	if err := ds.db.WithContext(ctx).Model(personnel).Updates(personnel).Error; err != nil {
		return nil, err
	}

	return personnel, nil
}

func (ds *DataStorage) DeletePersonnel(ctx context.Context, id string) error {
	personnel := &models.Personnel{
		Id: id,
	}
	if err := ds.db.WithContext(ctx).Where(personnel).Delete(personnel).Error; err != nil {
		return err
	}

	return nil
}

func (ds *DataStorage) GetPersonnel(ctx context.Context, id string) (*models.Personnel, error) {
	personnel := &models.Personnel{Id: id}
	err := ds.db.WithContext(ctx).Model(personnel).Take(personnel).Error
	if err != nil {
		if derrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return personnel, nil
}

// TODO can add sorting here
func (ds *DataStorage) ListPersonnels(ctx context.Context, page int, size int) ([]*models.Personnel, int, error) {
	res := []*models.Personnel{}
	smt := ds.db.WithContext(ctx).Model(&models.Personnel{}).Order("created_at ASC")
	if page > 0 && size > 0 {
		smt.Limit(size).Offset((page - 1) * size)
	}
	if err := smt.Find(&res).Error; err != nil {
		return nil, 0, err
	}
	count, err := ds.GetTotalCount(ctx, (&models.Personnel{}).TableName())
	return res, int(count), err
}
