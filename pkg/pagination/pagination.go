package pagination

type Pagination struct {
	Pn int `json:"page,omitempty" form:"page"`
	Ps int `json:"pageSize,omitempty" form:"pageSize"`
}

type Pager[T any] struct {
	Pagination
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}

func WrapPager[T any](total int64, list []T) *Pager[T] {
	return &Pager[T]{
		Total: total,
		List:  list,
	}
}
