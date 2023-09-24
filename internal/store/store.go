package store

type AppStore struct {
	userRepository UserRepository
	listRepository ListRepository
}

func NewAppStore(ur UserRepository, lr ListRepository) *AppStore {
	return &AppStore{
		userRepository: ur,
		listRepository: lr,
	}
}

func (s *AppStore) User() UserRepository {
	return s.userRepository
}

func (s *AppStore) List() ListRepository {
	return s.listRepository
}
