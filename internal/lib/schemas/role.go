package schemas

type RoleCreate struct {
	Name string `json:"role" binding:"required"`
}

type RoleUpdate struct {
	Name *string `json:"role" binding:"omitempty"`
}

type RoleRead struct {
	ID   int64  `json:"id"`
	Name string `json:"role"`
}
