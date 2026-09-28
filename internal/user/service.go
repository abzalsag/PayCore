package user

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateUser(email string, passwordHash string) (*User, error) {
	user := &User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	err := s.repository.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUserByID(id uint) (*User, error) {
	return s.repository.GetByID(id)
}
