package base

import "gorm.io/gorm"

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PageData[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}

func Wrapper[T any](total int64, list []T) *PageData[T] {
	return &PageData[T]{
		Total: total,
		List:  list,
	}
}

func Paginate(page Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.Page <= 0 {
			page.Page = 1
		}
		if page.PageSize <= 0 {
			page.PageSize = 10
		}
		offset := (page.Page - 1) * page.PageSize
		return db.Offset(offset).Limit(page.PageSize)
	}
}
