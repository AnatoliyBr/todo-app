package sqlstore

import (
	"database/sql"

	"github.com/AnatoliyBr/todo-app/internal/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(u *entity.User) error {
	return r.db.QueryRow("INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING user_id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.UserID)
}

func (r *UserRepository) FindByID(id int) (*entity.User, error) {
	return nil, nil
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	return nil, nil
}
