package storage

import (
	"context"
)

func (ds *DataStorage) GetTotalCount(ctx context.Context, tableName string) (int64, error) {
	var count int64
	err := ds.db.Table(tableName).WithContext(ctx).Count(&count).Error
	return count, err
}
