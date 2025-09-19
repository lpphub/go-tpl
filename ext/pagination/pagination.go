package pagination

type Pagination struct {
	Pn int `json:"page,omitempty" form:"page"`
	Ps int `json:"pageSize,omitempty" form:"pageSize"`
}

type PageData[T any] struct {
	Pagination
	Total int64 `json:"total"`
	List  []T   `json:"list,omitempty"`
}

func WrapPageData[T any](total int64, list []T) *PageData[T] {
	return &PageData[T]{
		Total: total,
		List:  list,
	}
}
