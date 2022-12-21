package db

import "gorm.io/gorm"

func Paging(db *gorm.DB, page, limit int) *gorm.DB {
	return db.Offset(limit * (page - 1)).Limit(limit)
}
