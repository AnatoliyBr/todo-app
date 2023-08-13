package sqlrepository

import (
	"database/sql"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/store"
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
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.db.QueryRow(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2) RETURNING user_id",
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.UserID)
}

func (r *UserRepository) FindByID(id int) (*entity.User, error) {
	u := &entity.User{}
	if err := r.db.QueryRow(
		"SELECT user_id, email, encrypted_password FROM users WHERE user_id = $1",
		id,
	).Scan(
		&u.UserID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	u := &entity.User{}
	if err := r.db.QueryRow(
		"SELECT user_id, email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.UserID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}
