package store

import "mvc/lego"

type Storage struct {
	models map[string]lego.Model
	names  []string
}

func NewStorage() Storage {
	return Storage{
		models: map[string]lego.Model{},
		names:  []string{},
	}
}

func (s *Storage) AddModel(model lego.Model) {
	s.models[model.Name] = model
	s.names = append(s.names, model.Name)
}

func (s *Storage) GetModelNames() []string {
	return s.names
}

func (s *Storage) GetModel(name string) lego.Model {
	return s.models[name]
}
