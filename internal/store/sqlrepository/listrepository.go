package sqlrepository

import (
	"database/sql"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/store"
)

type ListRepository struct {
	db *sql.DB
}

func NewListRepository(db *sql.DB) *ListRepository {
	return &ListRepository{
		db: db,
	}
}

func (r *ListRepository) Create(l *entity.List) error {
	if err := l.Validate(); err != nil {
		return err
	}

	return r.db.QueryRow(
		"INSERT INTO lists (list_title, user_id) VALUES ($1, $2) RETURNING list_id",
		l.ListTitle,
		l.UserID,
	).Scan(&l.ListID)
}

func (r *ListRepository) FindByID(listID, userID int) (*entity.List, error) {
	l := &entity.List{}
	if err := r.db.QueryRow(
		"SELECT list_id, list_title, user_id FROM lists WHERE list_id = $1 AND user_id = $2",
		listID,
		userID,
	).Scan(
		&l.ListID,
		&l.ListTitle,
		&l.UserID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return l, nil
}

func (r *ListRepository) Edit(l *entity.List) (*entity.List, error) {
	if err := l.Validate(); err != nil {
		return nil, err
	}

	_, err := r.db.Exec(
		"UPDATE lists SET list_title = $1 WHERE list_id = $2",
		l.ListTitle,
		l.ListID,
	)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *ListRepository) Delete(l *entity.List) error {
	_, err := r.db.Exec(
		"DELETE FROM lists WHERE list_id = $1",
		l.ListID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ListRepository) FindByUser(userID int) ([]*entity.List, error) {
	lists := make([]*entity.List, 0)

	rows, err := r.db.Query(
		"SELECT list_id, list_title FROM lists WHERE user_id = $1",
		userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var listID int
		var listTitle string

		err := rows.Scan(&listID, &listTitle)
		if err != nil {
			return nil, err
		}
		lists = append(lists, &entity.List{ListID: listID, ListTitle: listTitle})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(lists) > 0 {
		return lists, nil
	} else {
		return nil, store.ErrRecordNotFound
	}
}
