package role

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (Role) TableName() string {
	return "roles"
}
