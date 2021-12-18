package inmemory

import (
	"errors"
	"log"

	uuid "github.com/satori/go.uuid"
)

type Schema interface {
	Add(id uuid.UUID, value interface{})
	Delete(id uuid.UUID) error
	DeleteAll()
	Get(id uuid.UUID) interface{}
	Exists(id uuid.UUID) bool
}

type schema struct {
	values map[uuid.UUID]interface{}
}

func NewSchema() Schema {
	s := schema{}
	s.values = make(map[uuid.UUID]interface{})
	return &s
}

func (s *schema) Add(id uuid.UUID, value interface{}) {
	s.values[id] = value
	log.Println("add value in schema:", value)
}

func (s *schema) Delete(id uuid.UUID) error {
	if s.values == nil || len(s.values) <= 0 {
		return errors.New("error: no value exists")
	}

	if id == uuid.Nil || len(id) <= 0 {
		return errors.New("error: id not valid")
	}

	if !s.Exists(id) {
		return errors.New("error: value id not exists")
	}

	delete(s.values, id)
	return nil
}

func (s *schema) DeleteAll() {
	for k := range s.values {
		delete(s.values, k)
	}
}

func (s *schema) Exists(id uuid.UUID) bool {
	_, ok := s.values[id]
	return ok
}

func (s *schema) Get(id uuid.UUID) interface{} {
	return s.values[id]
}
