package schemas

type UserCreate struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserUpdate struct {
	Email    *string   `json:"email" binding:"omitempty,email"`
	Password *string   `json:"password" binding:"omitempty"`
	Roles    *[]string `json:"roles" binding:"omitempty,dive,required"`
}

type UserRead struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type UserReadByID struct {
	ID    string   `json:"id"`
	Email string   `json:"email"`
	Role  []string `json:"role"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
