package models

import "gorm.io/gorm"

func Pointer[T any](v T) *T {
	return &v
}

func Paginate(db *gorm.DB, page, pageSize int64) *gorm.DB {
	if page == -1 {
		return db
	}
	if page == 0 {
		page = 1
	}
	switch {
	case pageSize > 500:
		pageSize = 500
	case pageSize <= 0:
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return db.Offset(int(offset)).Limit(int(pageSize))
}
