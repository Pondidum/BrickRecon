package store

import "mvc/lego"

type Storage struct {
	models map[string]lego.Project
	names  []string
}

func NewStorage() Storage {
	return Storage{
		models: map[string]lego.Project{},
		names:  []string{},
	}
}

func (s *Storage) AddModel(model lego.Project) {
	s.models[model.Name] = model
	s.names = append(s.names, model.Name)
}

func (s *Storage) GetModelNames() []string {
	return s.names
}

func (s *Storage) GetModel(name string) lego.Project {
	return s.models[name]
}
