package postgres

import (
	"database/sql"
	"goP2Pbackend/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `INSERT INTO users (id, email, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Name, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) GetByID(id string) (*domain.User, error) {
	query := `SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1`
	var user domain.User
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, name, created_at, updated_at FROM users WHERE email = $1`
	var user domain.User
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	query := `UPDATE users SET email = $2, name = $3, updated_at = $4 WHERE id = $1`
	_, err := r.db.Exec(query, user.ID, user.Email, user.Name, user.UpdatedAt)
	return err
}
