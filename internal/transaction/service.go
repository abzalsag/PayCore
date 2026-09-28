package transaction

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetByID(id uint) (*Transaction, error) {
	return s.repository.GetByID(id)
}

func (s *Service) GetByWalletID(walletID uint) ([]Transaction, error) {
	return s.repository.GetByWalletID(walletID)
}
