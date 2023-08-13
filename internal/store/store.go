package store

type AppStore struct {
	userRepository UserRepository
}

func NewAppStore(ur UserRepository) *AppStore {
	return &AppStore{
		userRepository: ur,
	}
}

func (s *AppStore) User() UserRepository {
	return s.userRepository
}
