package user

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]User, error) {
	return s.repo.GetAll()
}

func (s *Service) GetByID(id int) (*User, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Create(u User) error {
	return s.repo.Create(u)
}

func (s *Service) UpdateFull(u UserUpdate) (int64, error) {
	return s.repo.UpdateFull(u)
}

func (s *Service) UpdatePartial(u UserUpdate) (int64, error) {
	return s.repo.UpdatePartial(u)
}

func (s *Service) Delete(id int) (int64, error) {
	return s.repo.Delete(id)
}
