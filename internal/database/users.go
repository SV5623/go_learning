package database

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (m *UserModel) Insert(user User) error {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return m.DB.QueryRow(query, user.Name, user.Email, user.Password).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)
}

func (m *UserModel) GetUser(query string, args ...interface{}) (*User, error) {
	ctx , cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	return m.GetUser(query, id)
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	return m.GetUser(query, email)
}