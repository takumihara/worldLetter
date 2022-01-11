package repository

import (
	"errors"
	"github.com/tacomea/worldLetter/domain"
	"sync"
)

type sessionRepository struct {
	m sync.Map
}

func NewSyncMapSessionRepository() domain.SessionRepository {
	return &sessionRepository{}
}

func (s *sessionRepository) Create(session domain.Session) error {
	s.m.Store(session.ID, session)
	return nil
}

func (s *sessionRepository) Delete(id string) error {
	s.m.Delete(id)
	return nil
}

func (s *sessionRepository) Read(id string) (domain.Session, error) {
	if value, ok := s.m.Load(id); ok {
		return value.(domain.Session), nil
	}
	return domain.Session{}, errors.New("user not found")
}
