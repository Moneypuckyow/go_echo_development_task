package user

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	DepartmentID int    `json:"department_id" validate:"required"`
}
