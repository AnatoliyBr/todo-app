package testrepository

import (
	"errors"

	"github.com/AnatoliyBr/todo-app/internal/entity"
	"github.com/AnatoliyBr/todo-app/internal/store"
)

type ListRepository struct {
	lists map[int]*entity.List
}

func NewListRepository() *ListRepository {
	return &ListRepository{
		lists: make(map[int]*entity.List),
	}
}

func (r *ListRepository) Create(l *entity.List) error {
	if err := l.Validate(); err != nil {
		return err
	}

	l.ListID = len(r.lists) + 1
	r.lists[l.ListID] = l

	return nil
}

func (r *ListRepository) FindByTitle(userID int, title string) (*entity.List, error) {
	for _, l := range r.lists {
		if l.UserID == userID && l.ListTitle == title {
			return l, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

func (r *ListRepository) Edit(l *entity.List) (*entity.List, error) {
	if err := l.Validate(); err != nil {
		return nil, err
	}

	for _, list := range r.lists {
		if list.UserID == l.UserID && list.ListTitle == l.ListTitle {
			return nil, errors.New("another list with this title has already exist")
		}
	}

	r.lists[l.ListID] = l
	return l, nil
}

func (r *ListRepository) Delete(l *entity.List) error {
	if _, ok := r.lists[l.ListID]; !ok {
		return store.ErrRecordNotFound
	}

	delete(r.lists, l.ListID)
	return nil
}

func (r *ListRepository) FindByUser(userID int) ([]*entity.List, error) {
	lists := make([]*entity.List, 0)

	for _, l := range r.lists {
		if l.UserID == userID {
			lists = append(lists, l)
		}
	}

	if len(lists) > 0 {
		return lists, nil
	} else {
		return nil, store.ErrRecordNotFound
	}
}
