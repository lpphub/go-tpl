package repo

import "gorm.io/gorm"

type EmptyList[T any] []T

func Empty[T any]() EmptyList[T] {
	return EmptyList[T]{}
}

func Paginate(pn, ps int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pn <= 0 {
			pn = 1
		}
		if ps <= 0 || ps > 200 {
			ps = 10
		}
		offset := (pn - 1) * ps
		return db.Offset(offset).Limit(ps)
	}
}
