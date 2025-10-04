package user

import (
	"database/sql"
	"fmt"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]User, error) {
	rows, err := r.db.Query(`SELECT user_id, name, email, department_id FROM "e-meeting"."user"`)
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
	row := r.db.QueryRow(`SELECT user_id, name, email, department_id FROM "e-meeting"."user" WHERE user_id = $1`, id)

	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.DepartmentID); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Create(u User) error {
	_, err := r.db.Exec(
		`INSERT INTO "e-meeting"."user" (name, email, department_id) VALUES ($1, $2, $3)`,
		u.Name, u.Email, u.DepartmentID,
	)
	return err
}

func (r *Repository) UpdateFull(u UserUpdate) (int64, error) {
	query := `UPDATE "e-meeting"."user" SET name = $1,email = $2,department_id = $3
	WHERE user_id = $4`
	res, err := r.db.Exec(query, u.Name, u.Email, u.DepartmentID, u.ID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (r *Repository) UpdatePartial(u UserUpdate) (int64, error) {
	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if u.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name=$%d", argPos))
		args = append(args, *u.Name)
		argPos++
	}
	if u.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email=$%d", argPos))
		args = append(args, *u.Email)
		argPos++
	}
	if u.DepartmentID != nil {
		setClauses = append(setClauses, fmt.Sprintf("department_id=$%d", argPos))
		args = append(args, *u.DepartmentID)
		argPos++
	}

	if len(setClauses) == 0 {
		return 0, nil // gak ada field yang diupdate
	}

	query := fmt.Sprintf(`UPDATE "e-meeting"."user" SET %s WHERE user_id=$%d`,
		strings.Join(setClauses, ", "), argPos)

	args = append(args, u.ID)

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *Repository) Delete(id int) (int64, error) {
	res, err := r.db.Exec(`DELETE FROM "e-meeting"."user" WHERE user_id=$1`, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
