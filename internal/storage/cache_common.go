package storage

import (
	"time"

	"example/pkg/cache"
	"example/pkg/log"
	"example/pkg/util"

	"gorm.io/gorm"
)

func (ds *DataStorage) DelCache(k string) error {
	if ds.cache != nil {
		return ds.cache.Del(k)
	}
	return nil
}

func (ds *DataStorage) Cache(k string, v interface{}, exp time.Duration) error {
	if ds.cache != nil {
		return ds.cache.SetStruct(k, v, exp)
	}
	return nil
}

func (ds *DataStorage) queryCache(out interface{}, keyCache string, exp time.Duration, dbExecute func(interface{}) error) error {
	if !util.IsStringEmpty(keyCache) && ds.cache != nil {
		if err := ds.cache.GetStruct(keyCache, out); err != nil {
			if err.Error() != cache.ErrCacheMiss.Error() {
				log.Errorw("Error when get data from cache", "error", err)
			}
		} else {
			return nil
		}
	}

	err := dbExecute(out)
	if err != nil {
		return err
	}

	// set back to cache
	if !util.IsStringEmpty(keyCache) && ds.cache != nil {
		if err := ds.cache.SetStruct(keyCache, out, exp); err != nil {
			log.Errorw("Error when set data from cache")
		}
	}
	return nil
}

// TakeCache caching Take query (get 1)
// WARN: be ware of this should be invalidated on update/delete
func (ds *DataStorage) TakeCache(db *gorm.DB, out interface{}, keyCache string, exp time.Duration) error {
	return ds.queryCache(out, keyCache, exp, func(o interface{}) error {
		return db.Take(o).Error
	})
}

// FindCache caching Find query (get all/limit)
// WARN: be ware of this should be invalidated on create/update/delete
// Should you on data that accept flushing async or rarely changing
func (ds *DataStorage) FindCache(db *gorm.DB, out interface{}, keyCache string, exp time.Duration) error {
	return ds.queryCache(out, keyCache, exp, func(o interface{}) error {
		return db.Find(o).Error
	})
}
