package storage

import (
	"context"
	"fmt"

	"example/internal/models"
	"example/pkg/cache"

	"go.uber.org/multierr"
	"gorm.io/gorm"
)

type DataStorage struct {
	db            *gorm.DB
	cache         cache.Cache
	inTransaction bool
}

func New(db *gorm.DB) *DataStorage {
	return &DataStorage{
		db:            db,
		inTransaction: false,
	}
}

func (ds *DataStorage) Close() error {
	db, err := ds.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (ds *DataStorage) SetCache(c cache.Cache) {
	ds.cache = c
}

func (ds *DataStorage) Migrate() error {
	return ds.db.AutoMigrate(
		&models.Personnel{}, &models.Education{})
}

func (ds *DataStorage) Transaction(ctx context.Context, f func(subds *DataStorage) error) (err error) {
	if ds.inTransaction {
		return f(ds)
	}
	transDB := ds.db.WithContext(ctx).Begin()
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("panic: %v", p)
			if rErr := transDB.Rollback().Error; rErr != nil {
				err = multierr.Combine(err, rErr)
			}
		}
	}()

	if err = f(&DataStorage{
		db:            transDB,
		inTransaction: true,
	}); err != nil {
		if rErr := transDB.Rollback().Error; rErr != nil {
			return multierr.Combine(err, rErr)
		}
		return err
	}

	err = transDB.Commit().Error
	return err
}
