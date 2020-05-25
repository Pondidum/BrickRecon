package store

type Storage struct {
	models []string
}

func NewStorage() Storage {
	return Storage{
		models: []string{"one", "two"},
	}
}

func (s *Storage) AddModel(name string) []string {

	s.models = append(s.models, name)

	return s.models
}

func (s *Storage) GetModels() []string {
	return s.models
}
