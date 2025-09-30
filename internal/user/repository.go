package user

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]User, error) {
	rows, err := r.db.Query(`SELECT user_id, name, email, department_id FROM "e-meeting"."User"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.DepartmentID); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) GetByID(id int) (*User, error) {
	row := r.db.QueryRow(`SELECT user_id, name, email, department_id FROM "e-meeting"."User" WHERE user_id = $1`, id)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.DepartmentID); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Create(u User) error {
	_, err := r.db.Exec(
		`INSERT INTO "e-meeting"."User" (name, email, department_id) VALUES ($1, $2, $3)`,
		u.Name, u.Email, u.DepartmentID,
	)
	return err
}

func (r *Repository) Update(u User) (int64, error) {
	res, err := r.db.Exec(
		`UPDATE "e-meeting"."User" SET name=$1, email=$2, department_id=$3 WHERE user_id=$4`,
		u.Name, u.Email, u.DepartmentID, u.ID,
	)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *Repository) Delete(id int) (int64, error) {
	res, err := r.db.Exec(`DELETE FROM "e-meeting"."User" WHERE user_id=$1`, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
