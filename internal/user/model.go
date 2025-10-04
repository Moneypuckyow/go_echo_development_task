package user

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	DepartmentID int    `json:"department_id" validate:"required"`
}

type UserUpdate struct {
	ID           int     `json:"id"`
	Name         *string `json:"name,omitempty"`
	Email        *string `json:"email,omitempty"`
	DepartmentID *int    `json:"department_id,omitempty"`
}
